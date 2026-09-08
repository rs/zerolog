package zerolog_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func consoleCountInput() []byte {
	var input bytes.Buffer
	logger := zerolog.New(&input)
	logger.Info().Str("foo", "bar").Msg("qux")
	return input.Bytes()
}

func TestConsoleWriterReturnsOriginalInputLength(t *testing.T) {
	input := consoleCountInput()
	original := append([]byte(nil), input...)
	var output bytes.Buffer
	writer := zerolog.ConsoleWriter{Out: &output, NoColor: true}
	n, err := writer.Write(input)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(input) {
		t.Errorf("Write returned %d bytes, want the %d input bytes", n, len(input))
	}
	if !bytes.Equal(input, original) {
		t.Error("Write modified its input")
	}
	if !strings.Contains(output.String(), "INF qux foo=bar") {
		t.Errorf("unexpected console output: %q", output.String())
	}
}

func TestConsoleWriterInIOMultiWriter(t *testing.T) {
	input := consoleCountInput()
	for _, consoleFirst := range []bool{false, true} {
		var raw, output bytes.Buffer
		console := zerolog.ConsoleWriter{Out: &output, NoColor: true}
		writers := []io.Writer{&raw, console}
		if consoleFirst {
			writers = []io.Writer{console, &raw}
		}
		n, err := io.MultiWriter(writers...).Write(input)
		if err != nil || n != len(input) {
			t.Errorf("console first=%v: Write = (%d, %v), want (%d, nil)", consoleFirst, n, err, len(input))
		}
		if !bytes.Equal(raw.Bytes(), input) {
			t.Errorf("console first=%v: raw writer did not receive the entire event", consoleFirst)
		}
	}
}

func TestConsoleWriterFormatterErrorsReturnZero(t *testing.T) {
	wantErr := errors.New("format failed")
	for _, prepare := range []bool{false, true} {
		var output bytes.Buffer
		writer := zerolog.ConsoleWriter{Out: &output, NoColor: true}
		if prepare {
			writer.FormatPrepare = func(map[string]interface{}) error { return wantErr }
		} else {
			writer.FormatExtra = func(map[string]interface{}, *bytes.Buffer) error { return wantErr }
		}
		n, err := writer.Write(consoleCountInput())
		if n != 0 || !errors.Is(err, wantErr) {
			t.Errorf("prepare=%v: Write = (%d, %v), want (0, %v)", prepare, n, err, wantErr)
		}
		if output.Len() != 0 {
			t.Errorf("prepare=%v: failed formatter wrote output", prepare)
		}
	}
}

type consoleCountErrorWriter struct{}

func (consoleCountErrorWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestConsoleWriterOutputErrorReportsOriginalInputLength(t *testing.T) {
	input := consoleCountInput()
	writer := zerolog.ConsoleWriter{Out: consoleCountErrorWriter{}, NoColor: true}
	n, err := writer.Write(input)
	if n != len(input) || !errors.Is(err, io.ErrClosedPipe) {
		t.Errorf("Write = (%d, %v), want (%d, %v)", n, err, len(input), io.ErrClosedPipe)
	}
}
