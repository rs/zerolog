package testdata

import (
	"github.com/rs/zerolog"
)

func p() {
	logger := zerolog.Nop()
	logger.Info().Int("key", 1) // want "missing or bad finisher for log chain, last call: \"Int\""
	logger.Info().Int("key", 1).Msg("abc")
}
