package zerolog_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestContextLoggerRef(t *testing.T) {
	out := &bytes.Buffer{}
	logger := zerolog.New(out)

	ctxD := logger.WithContext(context.Background())
	ctxR := logger.With().LoggerRef().WithContext(context.Background())

	// this is really less a test and more a demonstration that the issue
	// described in https://github.com/rs/zerolog/issues/116 is addressed
	assert.IsType(t, &logger, log.Ctx(ctxD))
	assert.IsType(t, &logger, log.Ctx(ctxR))
}
