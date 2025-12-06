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

func TestEvent_WithNilEvent(t *testing.T) {
	// coverage for nil Event receiver for all types
	var e *Event = nil

	fixtures := makeFieldFixtures()
	types := map[string]func() *Event{
		"Bool": func() *Event {
			return e.Bool("k", fixtures.Bools[0])
		},
		"Bools": func() *Event {
			return e.Bools("k", fixtures.Bools)
		},
		"Int": func() *Event {
			return e.Int("k", fixtures.Ints[0])
		},
		"Ints": func() *Event {
			return e.Ints("k", fixtures.Ints)
		},
		"Float": func() *Event {
			return e.Float64("k", fixtures.Floats[0])
		},
		"Floats": func() *Event {
			return e.Floats64("k", fixtures.Floats)
		},
		"Str": func() *Event {
			return e.Str("k", fixtures.Strings[0])
		},
		"Strs": func() *Event {
			return e.Strs("k", fixtures.Strings)
		},
		"Err": func() *Event {
			return e.Err(fixtures.Errs[0])
		},
		"Errs": func() *Event {
			return e.Errs("k", fixtures.Errs)
		},
		"Ctx": func() *Event {
			return e.Ctx(fixtures.Ctx)
		},
		"Time": func() *Event {
			return e.Time("k", fixtures.Times[0])
		},
		"Times": func() *Event {
			return e.Times("k", fixtures.Times)
		},
		"Dur": func() *Event {
			return e.Dur("k", fixtures.Durations[0])
		},
		"Durs": func() *Event {
			return e.Durs("k", fixtures.Durations)
		},
		"Interface": func() *Event {
			return e.Interface("k", fixtures.Interfaces[0])
		},
		"Interfaces": func() *Event {
			return e.Interface("k", fixtures.Interfaces)
		},
		"Interface(Object)": func() *Event {
			return e.Interface("k", fixtures.Objects[0])
		},
		"Interface(Objects)": func() *Event {
			return e.Interface("k", fixtures.Objects)
		},
		"Object": func() *Event {
			return e.Object("k", fixtures.Objects[0])
		},
		"Timestamp": func() *Event {
			return e.Timestamp()
		},
		"IPAddr": func() *Event {
			return e.IPAddr("k", fixtures.IPAddrs[0])
		},
		"IPAddrs": func() *Event {
			return e.IPAddrs("k", fixtures.IPAddrs)
		},
		"IPPrefix": func() *Event {
			return e.IPPrefix("k", fixtures.IPPfxs[0])
		},
		"IPPrefixes": func() *Event {
			return e.IPPrefixes("k", fixtures.IPPfxs)
		},
		"MACAddr": func() *Event {
			return e.MACAddr("k", fixtures.MACAddr)
		},
		"Type": func() *Event {
			return e.Type("k", fixtures.Type)
		},
	}

	for name := range types {
		f := types[name]
		if got := f(); got != nil {
			t.Errorf("Event.Bool() = %v, want %v", got, nil)
		}
	}
}
