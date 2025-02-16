package middleware

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func AppendToInterceptContext(ctx context.Context, key string, value []string) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	md = md.Copy()
	md[key] = value

	// Create a new context containing these modified metadata
	newCtx := metadata.NewIncomingContext(ctx, md)
	return newCtx
}
