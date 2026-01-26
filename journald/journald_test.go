//go:build linux
// +build linux

package journald_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/coreos/go-systemd/v22/journal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
)

func ExampleNewJournalDWriter() {
	log := zerolog.New(journald.NewJournalDWriter())
	log.Info().Str("foo", "bar").Uint64("small", 123).Float64("float", 3.14).Uint64("big", 1152921504606846976).Msg("Journal Test")
	// Output:
}

/*

There is no automated way to verify the output - since the output is sent
to journald process and method to retrieve is journalctl. Will find a way
to automate the process and fix this test.

$ journalctl -o verbose -f

Thu 2018-04-26 22:30:20.768136 PDT [s=3284d695bde946e4b5017c77a399237f;i=329f0;b=98c0dca0debc4b98a5b9534e910e7dd6;m=7a702e35dd4;t=56acdccd2ed0a;x=4690034cf0348614]
    PRIORITY=6
    _AUDIT_LOGINUID=1000
    _BOOT_ID=98c0dca0debc4b98a5b9534e910e7dd6
    _MACHINE_ID=926ed67eb4744580948de70fb474975e
    _HOSTNAME=sprint
    _UID=1000
    _GID=1000
    _CAP_EFFECTIVE=0
    _SYSTEMD_SLICE=-.slice
    _TRANSPORT=journal
    _SYSTEMD_CGROUP=/
    _AUDIT_SESSION=2945
    MESSAGE=Journal Test
    FOO=bar
    BIG=1152921504606846976
    _COMM=journald.test
    SMALL=123
    FLOAT=3.14
    JSON={"level":"info","foo":"bar","small":123,"float":3.14,"big":1152921504606846976,"message":"Journal Test"}
    _PID=27103
    _SOURCE_REALTIME_TIMESTAMP=1524807020768136
*/

func TestSanitizeKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", "TEST"},
		{"Test", "TEST"},
		{"test-key", "TEST_KEY"},
		{"Test.Key", "TEST_KEY"},
		{"test_key123", "TEST_KEY123"},
		{"invalid@key!", "INVALID_KEY_"},
		{"a1B2_c3D4", "A1B2_C3D4"},
		{"", ""},
		{"_", "_"},
		{"123", "123"},
		{"a-b.c_d", "A_B_C_D"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := journald.SanitizeKey(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeKey(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWriteReturnsNoOfWrittenBytes(t *testing.T) {
	input := []byte(`{"level":"info","time":1570912626,"message":"Starting..."}`)
	wr := journald.NewJournalDWriter()
	want := len(input)
	got, err := wr.Write(input)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if want != got {
		t.Errorf("Expected %d bytes to be written got %d", want, got)
	}
}

func TestMultiWrite(t *testing.T) {
	var (
		w1 = new(bytes.Buffer)
		w2 = new(bytes.Buffer)
		w3 = journald.NewJournalDWriter()
	)

	zerolog.ErrorHandler = func(err error) {
		if err == io.ErrShortWrite {
			t.Errorf("Unexpected ShortWriteError")
			t.FailNow()
		}
	}

	log := zerolog.New(io.MultiWriter(w1, w2, w3)).With().Logger()

	for i := 0; i < 10; i++ {
		log.Info().Msg("Tick!")
	}
}

func TestWriteWithVariousTypes(t *testing.T) {
	mock := &mockSend{}
	oldSend := journald.SendFunc
	journald.SendFunc = mock.send
	defer func() { journald.SendFunc = oldSend }()

	wr := journald.NewJournalDWriter()
	log := zerolog.New(wr)

	// This should cover the default case in the switch for value types
	log.Info().Bool("flag", true).Str("foo", "bar").Uint64("small", 123).Float64("float", 3.14).Uint64("big", 1152921504606846976).Interface("data", map[string]int{"a": 1}).Msg("Test various types")

	// Verify the call
	if len(mock.calls) != 1 {
		t.Fatalf("Expected 1 call, got %d", len(mock.calls))
	}

	call := mock.calls[0]

	// Check that flag is sanitized to FLAG and value is "true"
	if call.args["FLAG"] != "true" {
		t.Errorf("Expected FLAG=true, got %s", call.args["FLAG"])
	}

	// Check that data is marshaled (should be a JSON string)
	expectedData := `{"a":1}`
	if call.args["DATA"] != expectedData {
		t.Errorf("Expected DATA=%q, got %q", expectedData, call.args["DATA"])
	}
}

func TestWriteWithAllLevels(t *testing.T) {
	wr := journald.NewJournalDWriter()

	// Save original FatalExitFunc
	oldFatalExitFunc := zerolog.FatalExitFunc
	defer func() { zerolog.FatalExitFunc = oldFatalExitFunc }()

	// Set FatalExitFunc to prevent actual exit
	zerolog.FatalExitFunc = func() {}

	log := zerolog.New(wr)

	// Test all zerolog levels to cover levelToJPrio switch cases
	log.Trace().Msg("Trace level")
	log.Debug().Msg("Debug level")
	log.Info().Msg("Info level")
	log.Warn().Msg("Warn level")
	log.Error().Msg("Error level")
	log.Log().Msg("No level")

	// For Fatal, it will call FatalExitFunc instead of exiting
	log.Fatal().Msg("Fatal level")

	// For Panic, use recover to catch the panic, do last because it will stop of this test execution
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic from Panic level")
		}
	}()
	log.Panic().Msg("Panic level")

}

func TestWriteOutputs(t *testing.T) {
	mock := &mockSend{}
	oldSend := journald.SendFunc
	journald.SendFunc = mock.send
	defer func() { journald.SendFunc = oldSend }()

	wr := journald.NewJournalDWriter()
	log := zerolog.New(wr)

	// Log a message with various fields
	log.Info().Str("test-key", "value").Int("number", 42).Msg("Test message")

	// Check that SendFunc was called
	if len(mock.calls) != 1 {
		t.Fatalf("Expected 1 call to SendFunc, got %d", len(mock.calls))
	}

	call := mock.calls[0]

	// Check message
	if call.msg != "Test message" {
		t.Errorf("Expected msg 'Test message', got %q", call.msg)
	}

	// Check priority
	if call.prio != journal.PriInfo {
		t.Errorf("Expected prio %d (PriInfo), got %d", journal.PriInfo, call.prio)
	}

	// Check args
	expectedArgs := map[string]string{
		"TEST_KEY": "value",
		"NUMBER":   "42",
		"JSON":     `{"level":"info","test-key":"value","number":42,"message":"Test message"}` + "\n",
	}

	for k, v := range expectedArgs {
		if call.args[k] != v {
			t.Errorf("Expected args[%q] = %q, got %q", k, v, call.args[k])
		}
	}

	// Check that LEVEL is not in args (since it's skipped)
	if _, ok := call.args["LEVEL"]; ok {
		t.Error("LEVEL should not be in args")
	}
}

func TestWriteWithMarshalError(t *testing.T) {
	mock := &mockSend{}
	oldSend := journald.SendFunc
	journald.SendFunc = mock.send
	defer func() { journald.SendFunc = oldSend }()

	// Save original marshal func
	originalMarshal := zerolog.InterfaceMarshalFunc
	defer func() { zerolog.InterfaceMarshalFunc = originalMarshal }()

	// Set marshal func to fail
	zerolog.InterfaceMarshalFunc = func(v interface{}) ([]byte, error) {
		return nil, fmt.Errorf("fake error")
	}

	wr := journald.NewJournalDWriter()
	log := zerolog.New(wr)

	// This should trigger the error handling in the default case
	log.Info().Interface("data", map[string]int{"a": 1}).Msg("Test with error")

	// Verify the call
	if len(mock.calls) != 1 {
		t.Fatalf("Expected 1 call, got %d", len(mock.calls))
	}

	call := mock.calls[0]

	// Check that data has the error message
	got := call.args["DATA"]
	want := "error: fake error"
	if !strings.Contains(got, want) {
		t.Errorf("Expected DATA to contain %q, got %q", want, got)
	}
}

type mockSend struct {
	calls []struct {
		msg  string
		prio journal.Priority
		args map[string]string
	}
}

func (m *mockSend) send(msg string, prio journal.Priority, args map[string]string) error {
	m.calls = append(m.calls, struct {
		msg  string
		prio journal.Priority
		args map[string]string
	}{msg, prio, args})
	return nil
}
