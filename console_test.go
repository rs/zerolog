package zerolog_test

import (
	"os"

	"github.com/rs/zerolog"
)

func ExampleConsoleWriter_Write() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true})

	log.Info().Msg("hello world")
	// Output: <nil> |INFO| hello world
}
