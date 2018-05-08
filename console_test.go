package zerolog_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func ExampleConsoleWriter_Write() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true})

	log.Info().Msg("hello world")
	// Output: <nil> |INFO| hello world
}

func TestConsoleWriterNumbers(t *testing.T) {
	buf := &bytes.Buffer{}
	log := zerolog.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})
	log.Info().
		Float64("float", 1.23).
		Uint64("small", 123).
		Uint64("big", 1152921504606846976).
		Msg("msg")
	if got, want := strings.TrimSpace(buf.String()), "<nil> |INFO| msg big=1152921504606846976 float=1.23 small=123"; got != want {
		t.Errorf("\ngot:\n%s\nwant:\n%s", got, want)
	}
}
