package zerolog_test

import (
	"errors"
	stdlog "log"
	"os"
	"time"

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

func ExampleLogger_Error() {
	log := zerolog.New(os.Stdout)

	log.Error().
		Err(errors.New("some error")).
		Msg("error doing something")

	// Output: {"level":"error","error":"some error","message":"error doing something"}
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

func ExampleEvent_Durs() {
	d := []time.Duration{
		time.Duration(10 * time.Second),
		time.Duration(20 * time.Second),
	}

	log := zerolog.New(os.Stdout)

	log.Log().
		Str("foo", "bar").
		Durs("durs", d).
		Msg("hello world")

	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
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

func ExampleContext_Durs() {
	d := []time.Duration{
		time.Duration(10 * time.Second),
		time.Duration(20 * time.Second),
	}

	log := zerolog.New(os.Stdout).With().
		Str("foo", "bar").
		Durs("durs", d).
		Logger()

	log.Log().Msg("hello world")

	// Output: {"foo":"bar","durs":[10000,20000],"message":"hello world"}
}
