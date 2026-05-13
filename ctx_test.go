package zerolog

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/rs/zerolog/internal/cbor"
)

func TestCtx(t *testing.T) {
	log := New(io.Discard)
	ctx := log.WithContext(context.Background())
	log2 := Ctx(ctx)
	// Ctx attaches the Go context to the returned logger so events have it.
	want := log
	want.ctx = ctx
	if !reflect.DeepEqual(want, *log2) {
		t.Error("Ctx did not return the expected logger")
	}

	// update
	log = log.Level(InfoLevel)
	ctx = log.WithContext(ctx)
	log2 = Ctx(ctx)
	want = log
	want.ctx = ctx
	if !reflect.DeepEqual(want, *log2) {
		t.Error("Ctx did not return the expected logger")
	}

	log2 = Ctx(context.Background())
	if log2 != disabledLogger {
		t.Error("Ctx did not return the expected logger")
	}

	DefaultContextLogger = &log
	t.Cleanup(func() { DefaultContextLogger = nil })
	log2 = Ctx(context.Background())
	if log2 != &log {
		t.Error("Ctx did not return the expected logger")
	}
}

func TestCtxDisabled(t *testing.T) {
	dl := New(io.Discard).Level(Disabled)
	ctx := dl.WithContext(context.Background())
	if ctx != context.Background() {
		t.Error("WithContext stored a disabled logger")
	}

	l := New(io.Discard).With().Str("foo", "bar").Logger()
	ctx = l.WithContext(ctx)
	// Ctx returns a copy with ctx attached; compare against l with ctx set.
	lWithCtx := l
	lWithCtx.ctx = ctx
	if !reflect.DeepEqual(Ctx(ctx), &lWithCtx) {
		t.Error("WithContext did not store logger")
	}

	l.UpdateContext(func(c Context) Context {
		return c.Str("bar", "baz")
	})
	ctx = l.WithContext(ctx)
	lWithCtx = l
	lWithCtx.ctx = ctx
	if !reflect.DeepEqual(Ctx(ctx), &lWithCtx) {
		t.Error("WithContext did not store updated logger")
	}

	l = l.Level(DebugLevel)
	ctx = l.WithContext(ctx)
	lWithCtx = l
	lWithCtx.ctx = ctx
	if !reflect.DeepEqual(Ctx(ctx), &lWithCtx) {
		t.Error("WithContext did not store copied logger")
	}

	ctx = dl.WithContext(ctx)
	dlWithCtx := dl
	dlWithCtx.ctx = ctx
	if !reflect.DeepEqual(Ctx(ctx), &dlWithCtx) {
		t.Error("WithContext did not override logger with a disabled logger")
	}
}

type logObjectMarshalerImpl struct {
	name string
	age  int
}

func (t logObjectMarshalerImpl) MarshalZerologObject(e *Event) {
	e.Str("name", strings.ToLower(t.name)).Int("age", -t.age)
}

func Test_InterfaceLogObjectMarshaler(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf)
	ctx := log.WithContext(context.Background())

	log2 := Ctx(ctx)

	withLog := log2.With().Interface("obj", &logObjectMarshalerImpl{
		name: "FOO",
		age:  29,
	}).Logger()

	withLog.Info().Msg("test")

	if got, want := cbor.DecodeIfBinaryToString(buf.Bytes()), `{"level":"info","obj":{"name":"foo","age":-29},"message":"test"}`+"\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCtxEventHasContext(t *testing.T) {
	// Verify that events created from a context-retrieved logger automatically
	// have the Go context attached, so hooks can access it without callers
	// needing to chain .Ctx(ctx) on every event.
	var buf bytes.Buffer

	type ctxKeyType struct{}
	var gotCtx context.Context

	hook := HookFunc(func(e *Event, level Level, msg string) {
		gotCtx = e.GetCtx()
	})

	logger := New(&buf).Hook(hook)
	ctx := context.WithValue(context.Background(), ctxKeyType{}, "sentinel")
	ctx = logger.WithContext(ctx)

	Ctx(ctx).Info().Msg("test")

	if gotCtx == nil {
		t.Fatal("hook saw nil ctx; expected the context from WithContext")
	}
	if gotCtx.Value(ctxKeyType{}) != "sentinel" {
		t.Errorf("hook ctx missing expected value; got %v", gotCtx.Value(ctxKeyType{}))
	}
}
