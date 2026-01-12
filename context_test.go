package zerolog

import (
	"bytes"
	"errors"
	"testing"
)

type myError struct{}

func (e *myError) Error() string { return "test" }

func TestContext_ErrWithStackMarshaler(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler
	ErrorStackMarshaler = func(err error) interface{} {
		return "stack-trace"
	}

	var buf bytes.Buffer
	log := New(&buf).With().Stack().Err(errors.New("test error")).Logger()

	log.Info().Msg("test message")

	got := decodeIfBinaryToString(buf.Bytes())
	want := `{"level":"info","stack":"stack-trace","error":"test error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Context.Err() with stack marshaler = %q, want %q", got, want)
	}
}

func TestContext_AnErrWithNilErrorMarshal(t *testing.T) {
	// Save original
	original := ErrorMarshalFunc
	defer func() { ErrorMarshalFunc = original }()

	// Set marshaler to return a nil error pointer
	ErrorMarshalFunc = func(err error) interface{} {
		return (*myError)(nil) // nil pointer of error type
	}

	var buf bytes.Buffer
	log := New(&buf).With().AnErr("test", errors.New("some error")).Logger()

	log.Info().Msg("test message")

	got := decodeIfBinaryToString(buf.Bytes())
	want := `{"level":"info","message":"test message"}` + "\n" // No "test" field because isNilValue returned true
	if got != want {
		t.Errorf("Context.AnErr() with nil error marshal = %q, want %q", got, want)
	}
}

func TestContext_ErrWithNilStackMarshaler(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set marshaler to return nil
	ErrorStackMarshaler = func(err error) interface{} {
		return nil
	}

	var buf bytes.Buffer
	log := New(&buf).With().Stack().Err(errors.New("test error")).Logger()

	log.Info().Msg("test message")

	got := decodeIfBinaryToString(buf.Bytes())
	want := `{"level":"info","message":"test message"}` + "\n" // No stack or error field because stack marshaler returned nil
	if got != want {
		t.Errorf("Context.Err() with nil stack marshaler = %q, want %q", got, want)
	}
}

func TestContext_ErrWithStackMarshalerObject(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns LogObjectMarshaler
	ErrorStackMarshaler = func(err error) interface{} {
		return logObjectMarshalerImpl{name: "user", age: 30}
	}

	var buf bytes.Buffer
	log := New(&buf).With().Stack().Err(errors.New("test error")).Logger()

	log.Info().Msg("test message")

	got := decodeIfBinaryToString(buf.Bytes())
	want := `{"level":"info","stack":{"name":"user","age":-30},"error":"test error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Context.Err() with stack marshaler object = %q, want %q", got, want)
	}
}

func TestContext_ErrWithStackMarshalerError(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns an error
	ErrorStackMarshaler = func(err error) interface{} {
		return errors.New("stack error")
	}

	var buf bytes.Buffer
	log := New(&buf).With().Stack().Err(errors.New("test error")).Logger()

	log.Info().Msg("test message")

	got := decodeIfBinaryToString(buf.Bytes())
	want := `{"level":"info","stack":"stack error","error":"test error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Context.Err() with stack marshaler error = %q, want %q", got, want)
	}
}

func TestContext_ErrWithStackMarshalerInterface(t *testing.T) {
	// Save original
	original := ErrorStackMarshaler
	defer func() { ErrorStackMarshaler = original }()

	// Set a mock marshaler that returns an int
	ErrorStackMarshaler = func(err error) interface{} {
		return 42
	}

	var buf bytes.Buffer
	log := New(&buf).With().Stack().Err(errors.New("test error")).Logger()

	log.Info().Msg("test message")

	got := decodeIfBinaryToString(buf.Bytes())
	want := `{"level":"info","stack":42,"error":"test error","message":"test message"}` + "\n"
	if got != want {
		t.Errorf("Context.Err() with stack marshaler interface = %q, want %q", got, want)
	}
}
