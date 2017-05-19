package zerolog

import "context"

type ctxKey struct{}

// WithContext returns a copy of ctx with l associated.
func (l Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext returns the Logger associated with the ctx.
func FromContext(ctx context.Context) (l Logger, ok bool) {
	l, ok = ctx.Value(ctxKey{}).(Logger)
	return
}
