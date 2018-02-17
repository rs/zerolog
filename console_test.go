package zerolog_test

import (
	"os"

	"github.com/rs/zerolog"
)

func ExampleConsoleWriter_Write() {
	log := zerolog.New(zerolog.ConsoleWriter{os.Stdout, true})

	log.Info().Msg("hello world")
	// Output: <nil> |INFO| hello world
}
