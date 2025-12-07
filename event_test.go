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

func TestEvent_GetMetadata(t *testing.T) {
	type testCase struct {
		name    string
		e       *Event
		message string
		want    map[string]interface{}
	}

	testCases := []testCase{
		{
			name: "event without message",
			e:    newEvent(nil, DebugLevel).Str("foo", "bar").Float64("n", 42),
			want: map[string]interface{}{
				"foo": "bar",
				"n":   float64(42),
			},
		},
		{
			name: "event without message and integer",
			e:    newEvent(nil, DebugLevel).Str("foo", "bar").Int("n", 42),
			want: map[string]interface{}{
				"foo": "bar",
				"n":   float64(42),
			},
		},
		{
			name:    "event with message",
			e:       newEvent(nil, DebugLevel).Str("foo", "bar").Float64("n", 42),
			message: "test",
			want: map[string]interface{}{
				"foo":     "bar",
				"n":       float64(42),
				"message": "test",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.message != "" {
				tc.e.Msg(tc.message)
			}
			got, err := tc.e.GetMetadata()
			if err != nil {
				t.Error(err)
			}

			if len(got) != len(tc.want) {
				t.Errorf("Event.GetMetadata() = %v, want %v", len(got), len(tc.want))
			}
			for k, v := range tc.want {
				if got[k] != v {
					t.Errorf("Event.GetMetadata() = %v, want %v", got[k], v)
				}
			}
		})

	}
}
