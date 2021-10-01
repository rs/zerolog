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
