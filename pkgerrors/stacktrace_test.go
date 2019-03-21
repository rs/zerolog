// +build !binary_log

package pkgerrors

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func TestLogStack(t *testing.T) {
	zerolog.ErrorStackMarshaler = MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out)

	err := errors.Wrap(errors.New("error message"), "from error")
	log.Log().Stack().Err(err).Msg("")

	got := out.String()
	want := `\{"stack":\[\{"func":"TestLogStack","line":"20","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLogMultiStack(t *testing.T) {
	zerolog.ErrorStackMarshaler = MarshalMultiStack

	out := &bytes.Buffer{}
	log := zerolog.New(out)

	err := errors.Wrap(errors.New("error message"), "from error")
	log.Log().Stack().Err(err).Msg("")

	got := out.String()
	want := `\{"stack":\[\{"stacktrace":\[\{"source":"stacktrace_test.go","line":"36","func":"TestLogMultiStack"\},.*\{"stacktrace".*\],"error":"from error: error message"\}`
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}

}

// Some methods of wrapping cause more layers of wrapping than other layers,
// e.g. errors.New, errors.WithStack and errors.WithMessage add one layer of
// wrapping, whereas errors.Wrap adds two layers of wrapping.
func TestUnwrapErr(t *testing.T) {
	table := []struct {
		name               string
		err                error
		numberOfWrapLevels int
	}{
		{
			name:               "fundamental error",
			err:                errors.New("error message"),
			numberOfWrapLevels: 1,
		},
		{
			name:               "singly wrapped error",
			err:                errors.Wrap(errors.New("error message"), "from error"),
			numberOfWrapLevels: 3,
		},
		{
			name:               "doubly wrapped error",
			err:                errors.Wrap(errors.Wrap(errors.New("error message"), "first wrapper"), "second wrapper"),
			numberOfWrapLevels: 5,
		},
		{
			name:               "wrap with WithStack",
			err:                errors.WithStack(errors.New("error message")),
			numberOfWrapLevels: 2,
		},
		{
			name:               "wrap with WithMessage",
			err:                errors.WithMessage(errors.New("error message"), "first wrapper"),
			numberOfWrapLevels: 2,
		},
		{
			name:               "wrap with WithMessage and Wrap",
			err:                errors.Wrap(errors.WithMessage(errors.New("error message"), "first wrapper"), "second wrapper"),
			numberOfWrapLevels: 4,
		},
	}
	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			currentErr := test.err
			for i := 0; i < test.numberOfWrapLevels; i++ {
				currentErr = unwrapErr(currentErr)
			}
			if currentErr != nil {
				t.Fatal("Expected to have finished unwrapping by this point")
			}
		})
	}
}

func BenchmarkLogStack(b *testing.B) {
	zerolog.ErrorStackMarshaler = MarshalStack
	out := &bytes.Buffer{}
	log := zerolog.New(out)
	err := errors.Wrap(errors.New("error message"), "from error")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Log().Stack().Err(err).Msg("")
		out.Reset()
	}
}
