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

func ExampleBinaryConsoleWriter_Write() {
	log := zerolog.NewBinary(zerolog.ConsoleWriter{os.Stdout, true})

	log.Info().Int("Key", 100).Msg("hello world")
	// Output: <nil> |INFO| hello world Key=100
}
