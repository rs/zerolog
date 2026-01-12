//go:build !binary_log
// +build !binary_log

package zerolog

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

type nilError struct{}

func (nilError) Error() string {
	return "nope"
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
			e := newEvent(LevelWriterAdapter{&buf}, DebugLevel, false, nil, nil)
			e = e.AnErr("err", tt.err)
			err := e.write()
			if err != nil {
				t.Errorf("Event.AnErr() error: %v", err)
			}

			if got, want := strings.TrimSpace(buf.String()), tt.want; got != want {
				t.Errorf("Event.AnErr() = %v, want %v", got, want)
			}
		})
	}
}

func TestEvent_writeWithNil(t *testing.T) {
	var e *Event = nil
	got := e.write()

	var want *Event = nil
	if got != nil {
		t.Errorf("Event.write() = %v, want %v", got, want)
	}
}

type loggableObject struct {
	member string
}

func (o loggableObject) MarshalZerologObject(e *Event) {
	e.Str("member", o.member)
}

func TestEvent_Object(t *testing.T) {
	t.Run("ObjectWithNil", func(t *testing.T) {
		var buf bytes.Buffer
		e := newEvent(LevelWriterAdapter{&buf}, DebugLevel, false, nil, nil)
		e = e.Object("obj", nil)
		err := e.write()
		if err != nil {
			t.Errorf("Event.Object() error: %v", err)
		}

		want := `{"obj":null}`
		got := strings.TrimSpace(buf.String())
		if got != want {
			t.Errorf("Event.Object()\ngot:  %s\nwant: %s", got, want)
		}
	})

	t.Run("EmbedObjectWithNil", func(t *testing.T) {
		var buf bytes.Buffer
		e := newEvent(LevelWriterAdapter{&buf}, DebugLevel, false, nil, nil)
		e = e.EmbedObject(nil)
		err := e.write()
		if err != nil {
			t.Errorf("Event.EmbedObject() error: %v", err)
		}

		want := "{}"
		got := strings.TrimSpace(buf.String())
		if got != want {
			t.Errorf("Event.EmbedObject()\ngot:  %s\nwant: %s", got, want)
		}
	})

	type contextKeyType struct{}
	var contextKey = contextKeyType{}

	called := false
	ctxHook := HookFunc(func(e *Event, level Level, message string) {
		called = true
		ctx := e.GetCtx()
		if ctx == nil {
			t.Errorf("expected context to be set in Event")
		}
		val := ctx.Value(contextKey)
		if val == nil {
			t.Errorf("expected context value, got %v", val)
		}
		e.Str("ctxValue", val.(string))
		e.Bool("stackValue", e.stack)
	})

	t.Run("ObjectWithFullContext", func(t *testing.T) {
		called = false
		ctx := context.WithValue(context.Background(), contextKey, "ctx-object")

		var buf bytes.Buffer
		e := newEvent(LevelWriterAdapter{&buf}, DebugLevel, true, ctx, []Hook{ctxHook})
		e = e.Object("obj", loggableObject{member: "object-value"})
		e.Msg("hello")

		if !called {
			t.Errorf("hook was not called")
		}

		want := `{"obj":{"member":"object-value"},"ctxValue":"ctx-object","stackValue":true,"message":"hello"}`
		got := strings.TrimSpace(buf.String())
		if got != want {
			t.Errorf("Event.EmbedObject()\ngot:  %s\nwant: %s", got, want)
		}
	})

	t.Run("EmbedObjectWithFullContext", func(t *testing.T) {
		called = false
		ctx := context.WithValue(context.Background(), contextKey, "ctx-embed")

		var buf bytes.Buffer
		e := newEvent(LevelWriterAdapter{&buf}, DebugLevel, false, ctx, []Hook{ctxHook})
		e = e.EmbedObject(loggableObject{member: "embedded-value"})
		e.Msg("hello")

		if !called {
			t.Errorf("hook was not called")
		}

		want := `{"member":"embedded-value","ctxValue":"ctx-embed","stackValue":false,"message":"hello"}`
		got := strings.TrimSpace(buf.String())
		if got != want {
			t.Errorf("Event.EmbedObject()\ngot:  %s\nwant: %s", got, want)
		}
	})
}

func TestEvent_WithNilEvent(t *testing.T) {
	// coverage for nil Event receiver for all types
	var e *Event = nil

	fixtures := makeFieldFixtures()
	types := map[string]func() *Event{
		"Array": func() *Event {
			arr := e.CreateArray()
			return e.Array("k", arr)
		},
		"Bool": func() *Event {
			return e.Bool("k", fixtures.Bools[0])
		},
		"Bools": func() *Event {
			return e.Bools("k", fixtures.Bools)
		},
		"Fields": func() *Event {
			return e.Fields(fixtures)
		},
		"Int": func() *Event {
			return e.Int("k", fixtures.Ints[0])
		},
		"Ints": func() *Event {
			return e.Ints("k", fixtures.Ints)
		},
		"Int8": func() *Event {
			return e.Int8("k", fixtures.Ints8[0])
		},
		"Ints8": func() *Event {
			return e.Ints8("k", fixtures.Ints8)
		},
		"Int16": func() *Event {
			return e.Int16("k", fixtures.Ints16[0])
		},
		"Ints16": func() *Event {
			return e.Ints16("k", fixtures.Ints16)
		},
		"Int32": func() *Event {
			return e.Int32("k", fixtures.Ints32[0])
		},
		"Ints32": func() *Event {
			return e.Ints32("k", fixtures.Ints32)
		},
		"Int64": func() *Event {
			return e.Int64("k", fixtures.Ints64[0])
		},
		"Ints64": func() *Event {
			return e.Ints64("k", fixtures.Ints64)
		},
		"Uint": func() *Event {
			return e.Uint("k", fixtures.Uints[0])
		},
		"Uints": func() *Event {
			return e.Uints("k", fixtures.Uints)
		},
		"Uint8": func() *Event {
			return e.Uint8("k", fixtures.Uints8[0])
		},
		"Uints8": func() *Event {
			return e.Uints8("k", fixtures.Uints8)
		},
		"Uint16": func() *Event {
			return e.Uint16("k", fixtures.Uints16[0])
		},
		"Uints16": func() *Event {
			return e.Uints16("k", fixtures.Uints16)
		},
		"Uint32": func() *Event {
			return e.Uint32("k", fixtures.Uints32[0])
		},
		"Uints32": func() *Event {
			return e.Uints32("k", fixtures.Uints32)
		},
		"Uint64": func() *Event {
			return e.Uint64("k", fixtures.Uints64[0])
		},
		"Uints64": func() *Event {
			return e.Uints64("k", fixtures.Uints64)
		},
		"Float64": func() *Event {
			return e.Float64("k", fixtures.Floats64[0])
		},
		"Floats64": func() *Event {
			return e.Floats64("k", fixtures.Floats64)
		},
		"Float32": func() *Event {
			return e.Float32("k", fixtures.Floats32[0])
		},
		"Floats32": func() *Event {
			return e.Floats32("k", fixtures.Floats32)
		},
		"RawCBOR": func() *Event {
			return e.RawCBOR("k", fixtures.RawCBOR)
		},
		"RawJSON": func() *Event {
			return e.RawJSON("k", fixtures.RawJSONs[0])
		},
		"Str": func() *Event {
			return e.Str("k", fixtures.Strings[0])
		},
		"Strs": func() *Event {
			return e.Strs("k", fixtures.Strings)
		},
		"Stringers": func() *Event {
			return e.Stringers("k", fixtures.Stringers)
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
		"Dict": func() *Event {
			d := e.CreateDict()
			d.Str("greeting", "hello")
			return e.Dict("k", d)
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
		"Objects": func() *Event {
			return e.Objects("k", fixtures.Objects)
		},
		"EmbedObject": func() *Event {
			return e.EmbedObject(fixtures.Objects[0])
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
		"Caller": func() *Event {
			return e.Caller(1)
		},
		"CallerSkip": func() *Event {
			return e.CallerSkipFrame(2)
		},
		"Stack": func() *Event {
			return e.Stack()
		},
	}

	for name := range types {
		f := types[name]
		if got := f(); got != nil {
			t.Errorf("Event.Bool() = %v, want %v", got, nil)
		}
	}

	e.Send()
	e.Msg("nothing")
	e.Msgf("what %s", "nothing")

	got := e.write()
	if got != nil {
		t.Errorf("Event.write() = %v, want %v", got, e)
	}

	called := false
	e.MsgFunc(func() string {
		called = true
		return "called"
	})
	if called {
		t.Errorf("Event.MsgFunc() should not be called on nil Event")
	}
}

func TestEvent_MsgFunc(t *testing.T) {
	var buf bytes.Buffer
	e := newEvent(LevelWriterAdapter{&buf}, DebugLevel, false, nil, nil)

	called := false
	e.MsgFunc(func() string {
		called = true
		return "called"
	})
	if !called {
		t.Errorf("Event.MsgFunc() was not called on non-nil Event")
	}

	want := `{"message":"called"}`
	got := strings.TrimSpace(buf.String())
	if got != want {
		t.Errorf("Event.MsgFunc() = %q, want %q", got, want)
	}
}

func TestEvent_CallerRuntimeFail(t *testing.T) {
	var buf bytes.Buffer
	e := newEvent(LevelWriterAdapter{&buf}, DebugLevel, false, nil, nil)

	// Set a very large skipFrame to make runtime.Caller fail
	e.CallerSkipFrame(1000)
	e.Caller()

	e.Msg("test")

	got := strings.TrimSpace(buf.String())
	want := `{"message":"test"}` // No caller field because runtime.Caller failed
	if got != want {
		t.Errorf("Event.Caller() with failed runtime.Caller = %q, want %q", got, want)
	}
}

func TestEvent_DoneHandler(t *testing.T) {
	e := newEvent(nil, InfoLevel, false, nil, nil)

	// Set up a done handler to capture calls
	var called bool
	var capturedMsg string
	e.done = func(msg string) {
		called = true
		capturedMsg = msg
	}

	// Trigger msg via Msg
	e.Msg("test message")

	// Assert the handler was called with the correct message
	if !called {
		t.Error("Done handler was not called")
	}
	if capturedMsg != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", capturedMsg)
	}
}

type badLevelWriter struct {
	err error
}

func (w *badLevelWriter) WriteLevel(level Level, p []byte) (n int, err error) {
	return 0, w.err
}

func (w *badLevelWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}

func TestEvent_Msg_ErrorHandlerNil(t *testing.T) {
	// Save original ErrorHandler and restore after test
	originalErrorHandler := ErrorHandler
	ErrorHandler = nil
	defer func() { ErrorHandler = originalErrorHandler }()

	// Create a LevelWriter that always returns an error
	mockWriter := &badLevelWriter{err: errors.New("write error")}

	e := newEvent(mockWriter, InfoLevel, false, nil, nil)
	if e == nil {
		t.Fatal("Event should not be nil")
	}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	// Call Msg to trigger write error
	e.Msg("test message")

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = oldStderr
	captured, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	// Assert the error message was printed to stderr
	expected := "zerolog: could not write event: write error\n"
	if string(captured) != expected {
		t.Errorf("Expected stderr output %q, got %q", expected, string(captured))
	}
}

type mockLogObjectMarshaler struct {
	data string
}

func (m mockLogObjectMarshaler) MarshalZerologObject(e *Event) {
	e.Str("stack_func", m.data)
}

func TestEvent_ErrWithStackMarshaler(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler
	ErrorStackMarshaler = func(err error) interface{} {
		return "stack-trace"
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Err(err).Msg("test message")

	got := buf.String()
	want := `{"stack":"stack-trace","error":"test error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Event.Err() with stack marshaler = %q, want %q", got, want)
	}
}

func TestEvent_FieldsWithErrorAndStackMarshaler(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler
	ErrorStackMarshaler = func(err error) interface{} {
		return "stack-trace"
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Fields([]interface{}{"error", err}).Msg("test message")

	got := buf.String()
	want := `{"error":"test error","stack":"stack-trace","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Event.Fields() with error and stack marshaler = %q, want %q", got, want)
	}
}

func TestEvent_FieldsWithErrorAndStackMarshalerObject(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns LogObjectMarshaler
	ErrorStackMarshaler = func(err error) interface{} {
		return mockLogObjectMarshaler{data: "stack-data"}
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Fields([]interface{}{"error", err}).Msg("test message")

	got := buf.String()
	want := `{"error":"test error","stack":{"stack_func":"stack-data"},"message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Event.Fields() with error and stack marshaler object = %q, want %q", got, want)
	}
}

func TestEvent_FieldsWithErrorAndStackMarshalerError(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns an error
	ErrorStackMarshaler = func(err error) interface{} {
		return errors.New("stack error")
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Fields([]interface{}{"error", err}).Msg("test message")

	got := buf.String()
	want := `{"error":"test error","stack":"stack error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Event.Fields() with error and stack marshaler error = %q, want %q", got, want)
	}
}

func TestEvent_FieldsWithErrorAndStackMarshalerInterface(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns an int
	ErrorStackMarshaler = func(err error) interface{} {
		return 42
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Fields([]interface{}{"error", err}).Msg("test message")

	got := buf.String()
	want := `{"error":"test error","stack":42,"message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Event.Fields() with error and stack marshaler interface = %q, want %q", got, want)
	}
}

func TestEvent_FieldsWithErrorAndStackMarshalerNil(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set marshaler to return nil
	ErrorStackMarshaler = func(err error) interface{} {
		return nil
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Fields([]interface{}{"error", err}).Msg("test message")

	got := buf.String()
	want := `{"error":"test error","message":"test message"}` + "\n" // No stack field because marshaler returned nil
	if got != want {
		t.Errorf("Event.Fields() with error and nil stack marshaler = %q, want %q", got, want)
	}
}

func TestEvent_ErrWithStackMarshalerObject(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns LogObjectMarshaler
	ErrorStackMarshaler = func(err error) interface{} {
		return mockLogObjectMarshaler{data: "stack-data"}
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Err(err).Msg("test message")

	got := buf.String()
	want := `{"stack":{"stack_func":"stack-data"},"error":"test error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Event.Err() with stack marshaler object = %q, want %q", got, want)
	}
}

func TestEvent_ErrWithStackMarshalerError(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns an error
	ErrorStackMarshaler = func(err error) interface{} {
		return errors.New("stack error")
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Err(err).Msg("test message")

	got := buf.String()
	want := `{"stack":"stack error","error":"test error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Event.Err() with stack marshaler error = %q, want %q", got, want)
	}
}

func TestEvent_ErrWithStackMarshalerInterface(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns an int
	ErrorStackMarshaler = func(err error) interface{} {
		return 42
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Err(err).Msg("test message")

	got := buf.String()
	want := `{"stack":42,"error":"test error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Event.Err() with stack marshaler interface = %q, want %q", got, want)
	}
}

func TestEvent_ErrWithStackMarshalerNil(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set marshaler to return nil
	ErrorStackMarshaler = func(err error) interface{} {
		return nil
	}

	var buf bytes.Buffer
	log := New(&buf)

	err := errors.New("test error")
	log.Log().Stack().Err(err).Msg("test message")

	got := buf.String()
	want := `{"message":"test message"}` + "\n" // No fields because stack marshaler returned nil
	if got != want {
		t.Errorf("Event.Err() with nil stack marshaler = %q, want %q", got, want)
	}
}
