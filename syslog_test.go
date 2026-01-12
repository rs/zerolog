// +build !binary_log
// +build !windows

package zerolog

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

type syslogEvent struct {
	level string
	msg   string
}
type syslogTestWriter struct {
	events []syslogEvent
}

func (w *syslogTestWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
func (w *syslogTestWriter) Trace(m string) error {
	w.events = append(w.events, syslogEvent{"Trace", m})
	return nil
}
func (w *syslogTestWriter) Debug(m string) error {
	w.events = append(w.events, syslogEvent{"Debug", m})
	return nil
}
func (w *syslogTestWriter) Info(m string) error {
	w.events = append(w.events, syslogEvent{"Info", m})
	return nil
}
func (w *syslogTestWriter) Warning(m string) error {
	w.events = append(w.events, syslogEvent{"Warning", m})
	return nil
}
func (w *syslogTestWriter) Err(m string) error {
	w.events = append(w.events, syslogEvent{"Err", m})
	return nil
}
func (w *syslogTestWriter) Emerg(m string) error {
	w.events = append(w.events, syslogEvent{"Emerg", m})
	return nil
}
func (w *syslogTestWriter) Crit(m string) error {
	w.events = append(w.events, syslogEvent{"Crit", m})
	return nil
}

func TestSyslogWriter(t *testing.T) {
	sw := &syslogTestWriter{}
	log := New(SyslogLevelWriter(sw))
	log.Trace().Msg("trace")
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

type testCEEwriter struct {
	buf *bytes.Buffer
}

// Only implement one method as we're just testing the prefixing
func (c testCEEwriter) Debug(m string) error { return nil }

func (c testCEEwriter) Info(m string) error {
	_, err := c.buf.Write([]byte(m))
	return err
}

func (c testCEEwriter) Warning(m string) error { return nil }

func (c testCEEwriter) Err(m string) error { return nil }

func (c testCEEwriter) Emerg(m string) error { return nil }

func (c testCEEwriter) Crit(m string) error { return nil }

func (c testCEEwriter) Write(b []byte) (int, error) {
	return c.buf.Write(b)
}

func TestSyslogWriter_WithCEE(t *testing.T) {
	var buf bytes.Buffer
	sw := testCEEwriter{&buf}
	log := New(SyslogCEEWriter(sw))
	log.Info().Str("key", "value").Msg("message string")
	got := buf.String()
	want := "@cee:{"
	if !strings.HasPrefix(got, want) {
		t.Errorf("Bad CEE message start: want %v, got %v", want, got)
	}
}

type errorSyslogWriter struct {
	*syslogTestWriter
	writeError error
}

func (w *errorSyslogWriter) Write(p []byte) (int, error) {
	if w.writeError != nil {
		return 0, w.writeError
	}
	return len(p), nil
}

func TestSyslogWriter_Write(t *testing.T) {
	// Test Write method without prefix
	sw := &syslogTestWriter{}
	writer := SyslogLevelWriter(sw)

	data := []byte("test message")
	n, err := writer.Write(data)
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Write returned wrong length: got %d, want %d", n, len(data))
	}

	// Test Write method with CEE prefix
	sw2 := &syslogTestWriter{}
	writer2 := SyslogCEEWriter(sw2)

	data2 := []byte("test message")
	n2, err2 := writer2.Write(data2)
	if err2 != nil {
		t.Errorf("Write with CEE failed: %v", err2)
	}
	expectedLen := len(ceePrefix) + len(data2)
	if n2 != expectedLen {
		t.Errorf("Write with CEE returned wrong length: got %d, want %d", n2, expectedLen)
	}

	// Test Write method with CEE prefix and error on prefix write
	sw3 := &errorSyslogWriter{syslogTestWriter: &syslogTestWriter{}, writeError: io.EOF}
	writer3 := SyslogCEEWriter(sw3)

	_, err3 := writer3.Write(data2)
	if err3 != io.EOF {
		t.Errorf("Write with CEE error failed: got %v, want %v", err3, io.EOF)
	}
}

func TestSyslogWriter_WriteLevel_AllLevels(t *testing.T) {
	sw := &syslogTestWriter{}
	writer := SyslogLevelWriter(sw)

	// Test all levels to ensure full coverage
	writer.WriteLevel(TraceLevel, []byte(`{"level":"trace","message":"trace"}`+"\n"))
	writer.WriteLevel(DebugLevel, []byte(`{"level":"debug","message":"debug"}`+"\n"))
	writer.WriteLevel(InfoLevel, []byte(`{"level":"info","message":"info"}`+"\n"))
	writer.WriteLevel(WarnLevel, []byte(`{"level":"warn","message":"warn"}`+"\n"))
	writer.WriteLevel(ErrorLevel, []byte(`{"level":"error","message":"error"}`+"\n"))
	writer.WriteLevel(FatalLevel, []byte(`{"level":"fatal","message":"fatal"}`+"\n"))
	writer.WriteLevel(PanicLevel, []byte(`{"level":"panic","message":"panic"}`+"\n"))
	writer.WriteLevel(NoLevel, []byte(`{"message":"nolevel"}`+"\n"))

	want := []syslogEvent{
		{"Debug", `{"level":"debug","message":"debug"}` + "\n"},
		{"Info", `{"level":"info","message":"info"}` + "\n"},
		{"Warning", `{"level":"warn","message":"warn"}` + "\n"},
		{"Err", `{"level":"error","message":"error"}` + "\n"},
		{"Emerg", `{"level":"fatal","message":"fatal"}` + "\n"},
		{"Crit", `{"level":"panic","message":"panic"}` + "\n"},
		{"Info", `{"message":"nolevel"}` + "\n"},
	}
	if got := sw.events; !reflect.DeepEqual(got, want) {
		t.Errorf("Invalid syslog message routing: want %v, got %v", want, got)
	}
}

type closableSyslogWriter struct {
	*syslogTestWriter
	closed bool
}

func (w *closableSyslogWriter) Close() error {
	w.closed = true
	return nil
}

func TestSyslogWriter_Close(t *testing.T) {
	// Test with closable writer
	sw := &closableSyslogWriter{syslogTestWriter: &syslogTestWriter{}}
	writer := SyslogLevelWriter(sw).(syslogWriter) // Cast to concrete type to access Close

	err := writer.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
	if !sw.closed {
		t.Error("Close was not called on underlying writer")
	}

	// Test with non-closable writer
	sw2 := &syslogTestWriter{}
	writer2 := SyslogLevelWriter(sw2).(syslogWriter) // Cast to concrete type to access Close

	err = writer2.Close()
	if err != nil {
		t.Errorf("Close failed for non-closable writer: %v", err)
	}
}

func TestSyslogWriter_WriteLevel_InvalidLevel(t *testing.T) {
	sw := &syslogTestWriter{}
	writer := SyslogLevelWriter(sw)

	// Test invalid level - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid level")
		} else if r != "invalid level" {
			t.Errorf("Expected panic 'invalid level', got %v", r)
		}
	}()

	writer.WriteLevel(Level(100), []byte("test"))
}
