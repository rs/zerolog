// +build !binary_log
// +build !windows

package zerolog

import (
	"errors"
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
			name:    "All valid writers",
			writers: []io.Writer{
				mockedWriter {
					wantErr: false,
				},
				mockedWriter {
					wantErr: false,
				},
			},
		},
		{
			name:    "All invalid writers",
			writers: []io.Writer{
				mockedWriter {
					wantErr: true,
				},
				mockedWriter {
					wantErr: true,
				},
			},
		},
		{
			name:    "First invalid writer",
			writers: []io.Writer{
				mockedWriter {
					wantErr: true,
				},
				mockedWriter {
					wantErr: false,
				},
			},
		},
		{
			name:    "First valid writer",
			writers: []io.Writer{
				mockedWriter {
					wantErr: false,
				},
				mockedWriter {
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