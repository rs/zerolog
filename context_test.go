// +build !binary_log

package zerolog

import (
	"bytes"
	"context"
	"testing"
)

type logObjectMarshalerImpl struct {
	name string
	age  int
}

func (t logObjectMarshalerImpl) MarshalZerologObject(e *Event) {
	e.Str("name", "custom_value").Int("age", t.age)
}

func Test_InterfaceLogObjectMarshaler(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf)
	ctx := log.WithContext(context.Background())

	log2 := Ctx(ctx)

	withLog := log2.With().Interface("obj", &logObjectMarshalerImpl{
		name: "foo",
		age:  29,
	}).Logger()

	withLog.Info().Msg("test")

	if got, want := buf.String(), `{"level":"info","obj":{"name":"custom_value","age":29},"message":"test"}`+"\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
