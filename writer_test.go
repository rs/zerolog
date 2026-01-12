//go:build !binary_log && !windows
// +build !binary_log,!windows

package zerolog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"
)

type closableBuffer struct {
	*bytes.Buffer
	closed     bool
	closeError error
}

func (cb *closableBuffer) Close() error {
	cb.closed = true
	return cb.closeError
}

type errorWriter struct {
	writeError error
	shortWrite bool
}

func (ew *errorWriter) Write(p []byte) (int, error) {
	if ew.writeError != nil {
		return 0, ew.writeError
	}
	if ew.shortWrite {
		return len(p) - 1, nil // Return short write
	}
	return len(p), nil
}

func (ew *errorWriter) WriteLevel(level Level, p []byte) (int, error) {
	return ew.Write(p)
}

func TestMultiSyslogWriter(t *testing.T) {
	sw := &syslogTestWriter{}
	log := New(MultiLevelWriter(SyslogLevelWriter(sw)))
	log.Debug().Msg("debug")
	log.Info().Msg("info")
	log.Warn().Msg("warn")
	log.Error().Msg("error")
	log.Log().Msg("nolevel")
	want := []syslogEvent{
		{"Debug", `{"level":"debug","message":"debug"}` + "\n"},
		{"Info", `{"level":"info","message":"info"}` + "\n"},
		{"Warning", `{"level":"warn","message":"warn"}` + "\n"},
		{"Err", `{"level":"error","message":"error"}` + "\n"},
		{"Info", `{"message":"nolevel"}` + "\n"},
	}
	if got := sw.events; !reflect.DeepEqual(got, want) {
		t.Errorf("Invalid syslog message routing: want %v, got %v", want, got)
	}
}

var writeCalls int

type mockedWriter struct {
	wantErr bool
}

func (c mockedWriter) Write(p []byte) (int, error) {
	writeCalls++

	if c.wantErr {
		return -1, errors.New("Expected error")
	}

	return len(p), nil
}

// Tests that a new writer is only used if it actually works.
func TestResilientMultiWriter(t *testing.T) {
	tests := []struct {
		name    string
		writers []io.Writer
	}{
		{
			name: "All valid writers",
			writers: []io.Writer{
				mockedWriter{
					wantErr: false,
				},
				mockedWriter{
					wantErr: false,
				},
			},
		},
		{
			name: "All invalid writers",
			writers: []io.Writer{
				mockedWriter{
					wantErr: true,
				},
				mockedWriter{
					wantErr: true,
				},
			},
		},
		{
			name: "First invalid writer",
			writers: []io.Writer{
				mockedWriter{
					wantErr: true,
				},
				mockedWriter{
					wantErr: false,
				},
			},
		},
		{
			name: "First valid writer",
			writers: []io.Writer{
				mockedWriter{
					wantErr: false,
				},
				mockedWriter{
					wantErr: true,
				},
			},
		},
	}

	for _, tt := range tests {
		writers := tt.writers
		multiWriter := MultiLevelWriter(writers...)

		logger := New(multiWriter).With().Timestamp().Logger().Level(InfoLevel)
		logger.Info().Msg("Test msg")

		if len(writers) != writeCalls {
			t.Errorf("Expected %d writers to have been called but only %d were.", len(writers), writeCalls)
		}
		writeCalls = 0
	}
}

type testingLog struct {
	testing.TB
	buf bytes.Buffer
}

func (t *testingLog) Log(args ...interface{}) {
	if _, err := t.buf.WriteString(fmt.Sprint(args...)); err != nil {
		t.Error(err)
	}
}

func (t *testingLog) Logf(format string, args ...interface{}) {
	if _, err := t.buf.WriteString(fmt.Sprintf(format, args...)); err != nil {
		t.Error(err)
	}
}

func TestTestWriter(t *testing.T) {
	tests := []struct {
		name  string
		write []byte
		want  []byte
	}{{
		name:  "newline",
		write: []byte("newline\n"),
		want:  []byte("newline"),
	}, {
		name:  "oneline",
		write: []byte("oneline"),
		want:  []byte("oneline"),
	}, {
		name:  "twoline",
		write: []byte("twoline\n\n"),
		want:  []byte("twoline"),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &testingLog{TB: t} // Capture TB log buffer.
			w := TestWriter{T: tb}

			n, err := w.Write(tt.write)
			if err != nil {
				t.Error(err)
			}
			if n != len(tt.write) {
				t.Errorf("Expected %d write length but got %d", len(tt.write), n)
			}
			p := tb.buf.Bytes()
			if !bytes.Equal(tt.want, p) {
				t.Errorf("Expected %q, got %q.", tt.want, p)
			}

			log := New(NewConsoleWriter(ConsoleTestWriter(t)))
			log.Info().Str("name", tt.name).Msg("Success!")

			tb.buf.Reset()
		})
	}

}

func TestFilteredLevelWriter(t *testing.T) {
	buf := bytes.Buffer{}
	writer := FilteredLevelWriter{
		Writer: LevelWriterAdapter{&buf},
		Level:  InfoLevel,
	}
	_, err := writer.WriteLevel(DebugLevel, []byte("no"))
	if err != nil {
		t.Error(err)
	}
	_, err = writer.WriteLevel(InfoLevel, []byte("yes"))
	if err != nil {
		t.Error(err)
	}
	p := buf.Bytes()
	if want := "yes"; !bytes.Equal([]byte(want), p) {
		t.Errorf("Expected %q, got %q.", want, p)
	}
}

type testWrite struct {
	Level
	Line []byte
}

func TestTriggerLevelWriter(t *testing.T) {
	tests := []struct {
		write []testWrite
		want  []byte
		all   []byte
	}{{
		[]testWrite{
			{DebugLevel, []byte("no\n")},
			{InfoLevel, []byte("yes\n")},
		},
		[]byte("yes\n"),
		[]byte("yes\nno\n"),
	}, {
		[]testWrite{
			{DebugLevel, []byte("yes1\n")},
			{InfoLevel, []byte("yes2\n")},
			{ErrorLevel, []byte("yes3\n")},
			{DebugLevel, []byte("yes4\n")},
		},
		[]byte("yes2\nyes1\nyes3\nyes4\n"),
		[]byte("yes2\nyes1\nyes3\nyes4\n"),
	}}

	for k, tt := range tests {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			buf := bytes.Buffer{}
			writer := TriggerLevelWriter{Writer: LevelWriterAdapter{&buf}, ConditionalLevel: DebugLevel, TriggerLevel: ErrorLevel}
			t.Cleanup(func() { writer.Close() })
			for _, w := range tt.write {
				_, err := writer.WriteLevel(w.Level, w.Line)
				if err != nil {
					t.Error(err)
				}
			}
			p := buf.Bytes()
			if want := tt.want; !bytes.Equal([]byte(want), p) {
				t.Errorf("Expected %q, got %q.", want, p)
			}
			err := writer.Trigger()
			if err != nil {
				t.Error(err)
			}
			p = buf.Bytes()
			if want := tt.all; !bytes.Equal([]byte(want), p) {
				t.Errorf("Expected %q, got %q.", want, p)
			}
		})
	}
}

func TestLevelWriterAdapter_Close(t *testing.T) {
	// Test with closable writer
	buf := &bytes.Buffer{}
	adapter := LevelWriterAdapter{Writer: buf}

	// bytes.Buffer doesn't implement io.Closer, so Close should return nil
	err := adapter.Close()
	if err != nil {
		t.Errorf("Close should not return error for non-closable writer: %v", err)
	}

	// Test with closable writer
	closableBuf := &closableBuffer{Buffer: &bytes.Buffer{}}
	adapter2 := LevelWriterAdapter{Writer: closableBuf}

	err = adapter2.Close()
	if err != nil {
		t.Errorf("Close should not return error: %v", err)
	}
	if !closableBuf.closed {
		t.Error("Close should have been called on closable writer")
	}
}

func TestSyncWriter(t *testing.T) {
	buf := &bytes.Buffer{}

	// Test SyncWriter with regular io.Writer
	syncWriter := SyncWriter(buf)

	// Test Write
	data := []byte("test data")
	n, err := syncWriter.Write(data)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Write returned wrong length: got %d, want %d", n, len(data))
	}
	if got := buf.String(); got != string(data) {
		t.Errorf("Write wrote wrong data: got %q, want %q", got, string(data))
	}

	// Test SyncWriter with LevelWriter - use it with a logger
	levelBuf := &bytes.Buffer{}
	levelWriter := LevelWriterAdapter{levelBuf}
	syncLevelWriter := SyncWriter(levelWriter)

	logger := New(syncLevelWriter)
	logger.Info().Msg("test message")

	expected := `{"level":"info","message":"test message"}` + "\n"
	if got := levelBuf.String(); got != expected {
		t.Errorf("SyncWriter with LevelWriter failed: got %q, want %q", got, expected)
	}

	// Test SyncWriter Close with closable writer
	closableBuf := &closableBuffer{Buffer: &bytes.Buffer{}, closed: false}
	closableSyncWriter := SyncWriter(closableBuf)

	if closeable, ok := closableSyncWriter.(io.Closer); !ok {
		t.Error("SyncWriter should implement Close method")
	} else {
		err := closeable.Close()
		if err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}
	if !closableBuf.closed {
		t.Error("Close should have been called on closable writer")
	}

	// Test SyncWriter Close with closable writer that returns error
	errorBuf := &closableBuffer{Buffer: &bytes.Buffer{}, closed: false, closeError: io.EOF}
	errorSyncWriter := SyncWriter(errorBuf)

	if closeable, ok := errorSyncWriter.(io.Closer); !ok {
		t.Error("SyncWriter should implement Close method")
	} else {
		err := closeable.Close()
		if err != io.EOF {
			t.Errorf("Close should have returned EOF error, got: %v", err)
		}
	}
	if !errorBuf.closed {
		t.Error("Close should have been called on closable writer")
	}
}

func TestMultiLevelWriter_Write(t *testing.T) {
	// Test successful writes
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	multiWriter := MultiLevelWriter(buf1, buf2)

	data := []byte("test data")
	n, err := multiWriter.Write(data)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Write returned wrong length: got %d, want %d", n, len(data))
	}

	if got1 := buf1.String(); got1 != string(data) {
		t.Errorf("First writer got wrong data: got %q, want %q", got1, string(data))
	}
	if got2 := buf2.String(); got2 != string(data) {
		t.Errorf("Second writer got wrong data: got %q, want %q", got2, string(data))
	}

	// Test with error writer
	errorWriter1 := &errorWriter{writeError: io.EOF}
	buf3 := &bytes.Buffer{}

	errorMultiWriter := MultiLevelWriter(errorWriter1, buf3)

	_, err = errorMultiWriter.Write(data)
	if err != io.EOF {
		t.Errorf("Write should have returned EOF error, got: %v", err)
	}

	// Test with short write
	shortWriter := &errorWriter{shortWrite: true}
	buf4 := &bytes.Buffer{}

	shortMultiWriter := MultiLevelWriter(shortWriter, buf4)

	_, err = shortMultiWriter.Write(data)
	if err != io.ErrShortWrite {
		t.Errorf("Write should have returned ErrShortWrite, got: %v", err)
	}
}

func TestMultiLevelWriter_WriteLevel(t *testing.T) {
	// Test successful writes
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	multiWriter := MultiLevelWriter(buf1, buf2)

	data := []byte("test level data")
	n, err := multiWriter.WriteLevel(InfoLevel, data)
	if err != nil {
		t.Errorf("WriteLevel failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("WriteLevel returned wrong length: got %d, want %d", n, len(data))
	}

	if got1 := buf1.String(); got1 != string(data) {
		t.Errorf("First writer got wrong data: got %q, want %q", got1, string(data))
	}
	if got2 := buf2.String(); got2 != string(data) {
		t.Errorf("Second writer got wrong data: got %q, want %q", got2, string(data))
	}

	// Test with error writer
	errorWriter1 := &errorWriter{writeError: io.EOF}
	buf3 := &bytes.Buffer{}

	errorMultiWriter := MultiLevelWriter(errorWriter1, buf3)

	_, err = errorMultiWriter.WriteLevel(InfoLevel, data)
	if err != io.EOF {
		t.Errorf("WriteLevel should have returned EOF error, got: %v", err)
	}
}

func TestMultiLevelWriter_Close(t *testing.T) {
	buf1 := &closableBuffer{Buffer: &bytes.Buffer{}, closed: false}
	buf2 := &bytes.Buffer{} // non-closable

	multiWriter := MultiLevelWriter(buf1, buf2)

	// Cast to concrete type to access Close
	mw := multiWriter.(multiLevelWriter)

	err := mw.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if !buf1.closed {
		t.Error("First closable writer should have been closed")
	}

	// Test multiLevelWriter Close with error
	errorBuf1 := &closableBuffer{Buffer: &bytes.Buffer{}, closed: false, closeError: io.EOF}
	errorBuf2 := &bytes.Buffer{} // non-closable

	errorMultiWriter := MultiLevelWriter(errorBuf1, errorBuf2)
	emw := errorMultiWriter.(multiLevelWriter)

	err = emw.Close()
	if err != io.EOF {
		t.Errorf("Close should have returned EOF error, got: %v", err)
	}
	if !errorBuf1.closed {
		t.Error("First closable writer should have been closed")
	}
}

func TestNewTestWriter(t *testing.T) {
	writer := NewTestWriter(t)

	if writer.T != t {
		t.Error("NewTestWriter should set the testing interface")
	}
	if writer.Frame != 0 {
		t.Errorf("NewTestWriter should set Frame to 0, got %d", writer.Frame)
	}
}

func TestFilteredLevelWriter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	filteredWriter := FilteredLevelWriter{
		Writer: LevelWriterAdapter{buf},
		Level:  InfoLevel,
	}

	data := []byte("test data")
	n, err := filteredWriter.Write(data)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Write returned wrong length: got %d, want %d", n, len(data))
	}

	if got := buf.String(); got != string(data) {
		t.Errorf("Write should always write: got %q, want %q", got, string(data))
	}
}

func TestFilteredLevelWriter_Close(t *testing.T) {
	buf := &closableBuffer{Buffer: &bytes.Buffer{}, closed: false}
	filteredWriter := FilteredLevelWriter{
		Writer: LevelWriterAdapter{buf},
		Level:  InfoLevel,
	}

	err := filteredWriter.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if !buf.closed {
		t.Error("Underlying closable writer should have been closed")
	}

	// Test FilteredLevelWriter Close with error
	errorBuf := &closableBuffer{Buffer: &bytes.Buffer{}, closed: false, closeError: io.EOF}
	errorFilteredWriter := FilteredLevelWriter{
		Writer: LevelWriterAdapter{errorBuf},
		Level:  InfoLevel,
	}

	err = errorFilteredWriter.Close()
	if err != io.EOF {
		t.Errorf("Close should have returned EOF error, got: %v", err)
	}
	if !errorBuf.closed {
		t.Error("Underlying closable writer should have been closed")
	}
}
