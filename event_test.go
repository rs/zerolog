//go:build !binary_log
// +build !binary_log

package zerolog

import (
	"bytes"
	"errors"
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
			e := newEvent(LevelWriterAdapter{&buf}, DebugLevel)
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
	e := newEvent(LevelWriterAdapter{&buf}, DebugLevel)
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
	e := newEvent(LevelWriterAdapter{&buf}, DebugLevel)
	_ = e.EmbedObject(nil)
	_ = e.write()

	want := "{}"
	got := strings.TrimSpace(buf.String())
	if got != want {
		t.Errorf("Event.EmbedObject() = %q, want %q", got, want)
	}
}

func TestEvent_GetFields(t *testing.T) {
	e := newEvent(nil, DebugLevel)
	e.Str("foo", "bar").Float64("n", 42).Msg("test")

	got, err := e.GetFields()
	if err != nil {
		t.Error(err)
	}
	want := map[string]interface{}{
		"foo":     "bar",
		"n":       float64(42),
		"message": "test",
	}

	if len(got) != len(want) {
		t.Errorf("Event.GetFields() = %v, want %v", len(got), len(want))
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("Event.GetFields() = %v, want %v", got[k], v)
		}
	}
}
