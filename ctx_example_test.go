//go:build !binary_log

package zerolog_test

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

func ExampleLogger_WithContext() {
	parent := zerolog.New(os.Stdout).With().Str("component", "orders").Logger()
	ctx := parent.WithContext(context.Background())

	child := zerolog.Ctx(ctx).With().Str("order_id", "123").Logger()
	childCtx := child.WithContext(ctx)

	zerolog.Ctx(childCtx).Info().Msg("processing")
	zerolog.Ctx(ctx).Info().Msg("parent")

	// Output:
	// {"level":"info","component":"orders","order_id":"123","message":"processing"}
	// {"level":"info","component":"orders","message":"parent"}
}
