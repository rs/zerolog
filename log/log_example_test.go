//go:build !binary_log
// +build !binary_log

package log_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// setup would normally be an init() function, however, there seems
// to be something awry with the testing framework when we set the
// global Logger from an init()
func setup() {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = ""
	// In order to always output a static time to stdout for these
	// examples to pass, we need to override zerolog.TimestampFunc
	// and log.Logger globals -- you would not normally need to do this
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Simple logging example using the Print function in the log package
// Note that both Print and Printf are at the debug log level by default
func ExamplePrint() {
	setup()

	log.Print("hello world")
	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

// Simple logging example using the Printf function in the log package
func ExamplePrintf() {
	setup()

	log.Printf("hello %s", "world")
	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

// Example of a log with no particular "level"
func ExampleLog() {
	setup()
	log.Log().Msg("hello world")

	// Output: {"time":1199811905,"message":"hello world"}
}

// Example of a conditional level based on the presence of an error.
func ExampleErr() {
	setup()
	err := errors.New("some error")
	log.Err(err).Msg("hello world")
	log.Err(nil).Msg("hello world")

	// Output: {"level":"error","error":"some error","time":1199811905,"message":"hello world"}
	// {"level":"info","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "trace")
func ExampleTrace() {
	setup()
	log.Trace().Msg("hello world")

	// Output: {"level":"trace","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "debug")
func ExampleDebug() {
	setup()
	log.Debug().Msg("hello world")

	// Output: {"level":"debug","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "info")
func ExampleInfo() {
	setup()
	log.Info().Msg("hello world")

	// Output: {"level":"info","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "warn")
func ExampleWarn() {
	setup()
	log.Warn().Msg("hello world")

	// Output: {"level":"warn","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "error")
func ExampleError() {
	setup()
	log.Error().Msg("hello world")

	// Output: {"level":"error","time":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "fatal")
func ExampleFatal() {
	setup()
	err := errors.New("A repo man spends his life getting into tense situations")
	service := "myservice"

	log.Fatal().
		Err(err).
		Str("service", service).
		Msgf("Cannot start %s", service)

	// Outputs: {"level":"fatal","time":1199811905,"error":"A repo man spends his life getting into tense situations","service":"myservice","message":"Cannot start myservice"}
}

// Example of a log at a particular "level" (in this case, "panic")
func ExamplePanic() {
	setup()

	log.Panic().Msg("Cannot start")
	// Outputs: {"level":"panic","time":1199811905,"message":"Cannot start"} then panics
}

// This example uses command-line flags to demonstrate various outputs
// depending on the chosen log level.
func Example() {
	setup()
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Msg("This message appears only when log level set to Debug")
	log.Info().Msg("This message appears when log level set to Debug or Info")

	if e := log.Debug(); e.Enabled() {
		// Compute log output only if enabled.
		value := "bar"
		e.Str("foo", value).Msg("some debug message")
	}

	// Output: {"level":"info","time":1199811905,"message":"This message appears when log level set to Debug or Info"}
}

// Example of using the Output function in the log package to change the output destination
func ExampleOutput() {
	setup()

	out := &bytes.Buffer{}
	tee := log.Output(out)
	tee.Info().Msg("hello world")
	written := out.Len()

	log.Info().Int("bytes", written).Msg("wrote")
	// Output: {"level":"info","bytes":59,"time":1199811905,"message":"wrote"}
}

// Example of using the With function to add context fields
func ExampleWith() {
	setup()

	// you have to assign the result of With() to a new Logger and can't inline the level calls
	// because they need a *Logger receiver
	augmented := log.With().Str("service", "myservice").Logger()
	augmented.Info().Msg("hello world")
	// Output: {"level":"info","service":"myservice","time":1199811905,"message":"hello world"}
}

// Example of using the Level function to set the log level
func ExampleLevel() {
	setup()

	// you have to assign the result of Level() to a new Logger and can't inline the level calls
	// because they need a *Logger receiver
	leveled := log.Level(zerolog.ErrorLevel)
	leveled.Info().Msg("hello world")
	leveled.Error().Msg("I said HELLO")
	// Output: {"level":"error","time":1199811905,"message":"I said HELLO"}
}

type valueKeyType int

var valueKey valueKeyType = 42

var captainHook = zerolog.HookFunc(func(e *zerolog.Event, l zerolog.Level, msg string) {
	e.Interface("key", e.GetCtx().Value(valueKey))
	e.Bool("is_error", l > zerolog.ErrorLevel)
	e.Int("msg_len", len(msg))
})

// Example of using the Logger Hook function to add hooks
func ExampleLogger_Hook() {
	setup()

	hooked := log.Hook(captainHook)
	hooked.Info().Msg("watch out!")
	// Output: {"level":"info","time":1199811905,"key":null,"is_error":false,"msg_len":10,"message":"watch out!"}
}

// Example of using the WithLevel function to set the log level
func ExampleWithLevel() {
	setup()

	// you have to assign the result of Level() to a new Logger and can't inline the level calls
	// because they need a *Logger receiver
	event := log.WithLevel(zerolog.ErrorLevel)
	event.Msg("taxes are due")
	// Output: {"level":"error","time":1199811905,"message":"taxes are due"}
}

// Example of using the Ctx function in the log package to log with context
func ExampleCtx() {
	setup()

	hooked := log.Hook(captainHook)
	ctx := context.WithValue(context.Background(), valueKey, "12345")
	logger := hooked.With().Ctx(ctx).Logger()
	log.Ctx(logger.WithContext(ctx)).Info().Msg("hello world")
	// Output: {"level":"info","time":1199811905,"key":"12345","is_error":false,"msg_len":11,"message":"hello world"}
}

// Example of using the Sample function in the log package to set a sampler
func ExampleSample() {
	setup()

	sampled := log.Sample(&zerolog.BasicSampler{N: 2})
	sampled.Info().Msg("hello world")
	sampled.Info().Msg("I said, hello world")
	sampled.Info().Msg("Can you here me now world")
	// Output: {"level":"info","time":1199811905,"message":"hello world"}
	// {"level":"info","time":1199811905,"message":"Can you here me now world"}
}
