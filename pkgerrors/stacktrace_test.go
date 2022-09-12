// +build !binary_log

package pkgerrors_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func TestLogStack(t *testing.T) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out)

	err := pkgerrors.Wrap(pkgerrors.New("error message"), "from error")
	log.Log().Stack().Err(err).Msg("")

	got := out.String()
	want := `\{"stack":\[\{"func":"TestLogStack","line":"20","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLogStackFromContext(t *testing.T) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out).With().Stack().Logger() // calling Stack() on log context instead of event

	err := pkgerrors.Wrap(pkgerrors.New("error message"), "from error")
	log.Log().Err(err).Msg("") // not explicitly calling Stack()

	got := out.String()
	want := `\{"stack":\[\{"func":"TestLogStackFromContext","line":"36","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func BenchmarkLogStack(b *testing.B) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	out := &bytes.Buffer{}
	log := zerolog.New(out)
	err := pkgerrors.Wrap(pkgerrors.New("error message"), "from error")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Log().Stack().Err(err).Msg("")
		out.Reset()
	}
}
