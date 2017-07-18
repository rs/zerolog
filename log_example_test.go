package zerolog_test

import (
	stdErrors "errors"
	stdlog "log"
	"os"
	"time"

	errors "github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func ExampleNew() {
	log := zerolog.New(os.Stdout)

	log.Info().Msg("hello world")

	// Output: {"level":"info","message":"hello world"}
}

func ExampleLogger_With() {
	log := zerolog.New(os.Stdout).
		With().
		Str("foo", "bar").
		Logger()

	log.Info().Msg("hello world")

	// Output: {"level":"info","foo":"bar","message":"hello world"}
}

func ExampleLogger_Level() {
	log := zerolog.New(os.Stdout).Level(zerolog.WarnLevel)

	log.Info().Msg("filtered out message")
	log.Error().Msg("kept message")

	// Output: {"level":"error","message":"kept message"}
}

func ExampleLogger_Sample() {
	log := zerolog.New(os.Stdout).Sample(2)

	log.Info().Msg("message 1")
	log.Info().Msg("message 2")
	log.Info().Msg("message 3")
	log.Info().Msg("message 4")

	// Output: {"level":"info","sample":2,"message":"message 2"}
	// {"level":"info","sample":2,"message":"message 4"}
}

func ExampleLogger_Debug() {
	log := zerolog.New(os.Stdout)

	log.Debug().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	// Output: {"level":"debug","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Info() {
	log := zerolog.New(os.Stdout)

	log.Info().
		Str("foo", "bar").
		Int("n", 123).
		Msg("hello world")

	// Output: {"level":"info","foo":"bar","n":123,"message":"hello world"}
}

func ExampleLogger_Warn() {
	log := zerolog.New(os.Stdout)

	log.Warn().
		Str("foo", "bar").
		Msg("a warning message")

	// Output: {"level":"warn","foo":"bar","message":"a warning message"}
}

func f1() error {
	return errors.Wrap(stdErrors.New("my nested error"), "f1")
}

func ExampleLogger_Error() {
	log := zerolog.New(os.Stdout)

	log.Error().
		Err(stdErrors.New("my standard error")).
		Msg("error doing something standard")

	log.Error().
		Err(stdErrors.New("my standard error want stack")).
		StackTrace().
		Msg("error doing something standard want stack")

	log.Error().
		Error("error doing something error-ish (no stack)")

	log.Error().
		StackTrace().
		Error("error doing something error-ish (want stack)")

	log.Error().
		Err(errors.New("my wrapped error")).
		Msg("error doing something wrapped (no stack)")

	log.Error().
		Err(errors.New("my wrapped error")).
		StackTrace().
		Msg("error doing something wrapped (with stack)")

	log.Error().
		Err(f1()).
		Msg("error doing something nested and wrapped (no stack)")

	log.Error().
		Err(f1()).
		StackTrace().
		Msg("error doing something nested and wrapped (with stack)")

	// NOTE(seanc@): This test output is known to be brittle due to line-numbers
	// and source filenames.

	// Output: {"level":"error","message":"error doing something standard","error":"my standard error"}
	// {"level":"error","message":"error doing something standard want stack","stack":[{"file":"event.go","line":"136","func":"(*Event).Msg"},{"file":"log_example_test.go","line":"99","func":"ExampleLogger_Error"},{"file":"example.go","line":"122","func":"runExample"},{"file":"example.go","line":"46","func":"runExamples"},{"file":"testing.go","line":"823","func":"(*M).Run"},{"file":"_testmain.go","line":"140","func":"main"},{"file":"proc.go","line":"185","func":"main"},{"file":"asm_amd64.s","line":"2197","func":"goexit"}]}
	// {"level":"error","message":"error doing something error-ish (no stack)","error":"error doing something error-ish (no stack)"}
	// {"level":"error","message":"error doing something error-ish (want stack)","stack":[{"file":"event.go","line":"136","func":"(*Event).Msg"},{"file":"event.go","line":"83","func":"(*Event).Error"},{"file":"log_example_test.go","line":"106","func":"ExampleLogger_Error"},{"file":"example.go","line":"122","func":"runExample"},{"file":"example.go","line":"46","func":"runExamples"},{"file":"testing.go","line":"823","func":"(*M).Run"},{"file":"_testmain.go","line":"140","func":"main"},{"file":"proc.go","line":"185","func":"main"},{"file":"asm_amd64.s","line":"2197","func":"goexit"}]}
	// {"level":"error","message":"error doing something wrapped (no stack)","error":"my wrapped error"}
	// {"level":"error","message":"error doing something wrapped (with stack)","stack":[{"file":"event.go","line":"136","func":"(*Event).Msg"},{"file":"log_example_test.go","line":"115","func":"ExampleLogger_Error"},{"file":"example.go","line":"122","func":"runExample"},{"file":"example.go","line":"46","func":"runExamples"},{"file":"testing.go","line":"823","func":"(*M).Run"},{"file":"_testmain.go","line":"140","func":"main"},{"file":"proc.go","line":"185","func":"main"},{"file":"asm_amd64.s","line":"2197","func":"goexit"}]}
	// {"level":"error","message":"error doing something nested and wrapped (no stack)","error":"f1: my nested error"}
	// {"level":"error","message":"error doing something nested and wrapped (with stack)","error":"f1: my nested error","stack":[{"file":"log_example_test.go","line":"86","func":"f1"},{"file":"log_example_test.go","line":"123","func":"ExampleLogger_Error"},{"file":"example.go","line":"122","func":"runExample"},{"file":"example.go","line":"46","func":"runExamples"},{"file":"testing.go","line":"823","func":"(*M).Run"},{"file":"_testmain.go","line":"140","func":"main"},{"file":"proc.go","line":"185","func":"main"},{"file":"asm_amd64.s","line":"2197","func":"goexit"}]}
}

func ExampleLogger_WithLevel() {
	log := zerolog.New(os.Stdout)

	log.WithLevel(zerolog.InfoLevel).
		Msg("hello world")

	// Output: {"level":"info","message":"hello world"}
}

func ExampleLogger_Write() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(log)

	stdlog.Print("hello world")

	// Output: {"foo":"bar","message":"hello world"}
}

func ExampleLogger_Log() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Str("bar", "baz").
		Msg("")

	// Output: {"foo":"bar","bar":"baz"}
}

func ExampleEvent_Dict() {
	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).
		Msg("hello world")

	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
}

func ExampleEvent_Interface() {
	log := zerolog.New(os.Stdout)

	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log.Log().
		Str("foo", "bar").
		Interface("obj", obj).
		Msg("hello world")

	// Output: {"foo":"bar","obj":{"name":"john"},"message":"hello world"}
}

func ExampleEvent_Dur() {
	d := time.Duration(10 * time.Second)

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Dur("dur", d).
		Msg("hello world")

	// Output: {"foo":"bar","dur":10000,"message":"hello world"}
}

func ExampleContext_Dict() {
	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		).Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","dict":{"bar":"baz","n":1},"message":"hello world"}
}

func ExampleContext_Interface() {
	obj := struct {
		Name string `json:"name"`
	}{
		Name: "john",
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Interface("obj", obj).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","obj":{"name":"john"},"message":"hello world"}
}

func ExampleContext_Dur() {
	d := time.Duration(10 * time.Second)

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Dur("dur", d).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","dur":10000,"message":"hello world"}
}
