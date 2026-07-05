package logger

import (
	"context"
)

type loggerKey struct{}

func LoggerFromContext[T any](ctx context.Context, defaultLogger T) T {
	if l, ok := ctx.Value(loggerKey{}).(T); ok {
		return l
	}

	return *new(T)
}

func WithLogger[T any](ctx context.Context, logger T) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}
