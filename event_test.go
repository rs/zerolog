//go:build !binary_log
// +build !binary_log

package zerolog

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type nilError struct{}

func (nilError) Error() string {
	return ""
}

func TestEvent_AnErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil", nil, `{}`},
		{"error", errors.New("test"), `{"err":"test"}`},
		{"nil interface", func() *nilError { return nil }(), `{}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			e := newEvent(levelWriterAdapter{&buf}, DebugLevel)
			e.AnErr("err", tt.err)
			_ = e.write()
			if got, want := strings.TrimSpace(buf.String()), tt.want; got != want {
				t.Errorf("Event.AnErr() = %v, want %v", got, want)
			}
		})
	}
}

func TestEvent_ObjectWithNil(t *testing.T) {
	var buf bytes.Buffer
	e := newEvent(levelWriterAdapter{&buf}, DebugLevel)
	_ = e.Object("obj", nil)
	_ = e.write()

	want := `{"obj":null}`
	got := strings.TrimSpace(buf.String())
	if got != want {
		t.Errorf("Event.Object() = %q, want %q", got, want)
	}
}

func TestEvent_EmbedObjectWithNil(t *testing.T) {
	var buf bytes.Buffer
	e := newEvent(levelWriterAdapter{&buf}, DebugLevel)
	_ = e.EmbedObject(nil)
	_ = e.write()

	want := "{}"
	got := strings.TrimSpace(buf.String())
	if got != want {
		t.Errorf("Event.EmbedObject() = %q, want %q", got, want)
	}
}

func TestEvent_LogLevel(t *testing.T) {
	levels := []Level{
		TraceLevel,
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
		PanicLevel,
		NoLevel,
		Disabled,
	}

	for _, l := range levels {
		t.Run(l.String(), func(t *testing.T) {
			out := &bytes.Buffer{}
			log := New(out)
			log.WithLevel(l).LogLevel().Msg("test")
			if l == Disabled {
				if got := decodeIfBinaryToString(out.Bytes()); got != `` {
					t.Errorf("invalid log output:\ngot:  %v\nwant: ", got)
				}
				return
			}
			if got, want := decodeIfBinaryToString(out.Bytes()), fmt.Sprintf(`{"level":"%s","message":"test"}`+"\n", l); got != want {
				t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
			}
		})
	}
}
