package zerolog

import "context"

// Hook defines an interface to a contex hook.
type CtxtHook interface {
	// Run runs the hook with extracting of the Logger associated with the ctx.
	Run(ctx context.Context, l Logger) Logger
}

// CtxtHookFunc is an adaptor to allow the use of an ordinary function as a Hook.
type CtxtHookFunc func(ctx context.Context, l Logger) Logger

// Run implements the CtxtHook interface.
func (h CtxtHookFunc) Run(ctx context.Context, l Logger) Logger {
	return h(ctx, l)
}
