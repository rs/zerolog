package log_test

import (
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {

	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}

	//log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Simple logging example using the log package.
func ExamplePrint() {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = ""

	log.Print("hello world")
	// Output: {"time":1199811905,"level":"debug","message":"hello world"}
}

// Simple example of a log at a particular "level" (in this case, "info")
// using the log package.  You'll note the function ExampleNew below is similar,
// but has a different log output as it's getting an instance of
// the logger from the zerolog package instead of using the log package
func ExampleInfo() {
	zerolog.TimeFieldFormat = ""

	log.Info().Msg("hello world")
	// Output: {"time":1199811905,"level":"info","message":"hello world"}
}

func ExampleFatal() {
	err := errors.New("A repo man spends his life getting into tense situations")
	service := "myservice"

	zerolog.TimeFieldFormat = ""

	log.Fatal().
		Err(err).
		Str("service", service).
		Msgf("Cannot start %s", service)
	// Outputs: {"time":1199811905,"level":"fatal","error":"A repo man spends his life getting into tense situations","service":"myservice","message":"Cannot start myservice"}
}
