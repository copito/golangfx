package authbypass

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/pkg/middleware"
	"google.golang.org/grpc"
)

type LocalBypassAuthInterceptor struct {
	logger *slog.Logger
	config *entities.Config
}

func NewLocalBypassAuthInterceptor(logger *slog.Logger, config *entities.Config) *LocalBypassAuthInterceptor {
	return &LocalBypassAuthInterceptor{
		logger: logger,
		config: config,
	}
}

// BuildUnaryInterceptor is a unary interceptor that adds a auth token to context if running locally
func (i *LocalBypassAuthInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		i.logger.Debug("middleware: local bypass auth", slog.String("full_method", info.FullMethod))

		// Check if environment is local
		if i.config.Backend.Environment != "local" {
			return handler(ctx, req)
		}

		// Add header bearer token (fake token)
		const fakeToken string = "example_fake_token"
		ctx = middleware.AppendToInterceptContext(ctx, "authorization", []string{fmt.Sprintf("Bearer %s", fakeToken)})

		i.logger.Debug("adding fake token to authorization", slog.String("token", fakeToken))
		resp, err := handler(ctx, req)
		return resp, err
	}
}

// StreamInterceptor us a stream interceptor that adds a auth token to context if running locally
func (i *LocalBypassAuthInterceptor) BuildStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		i.logger.Debug("running StreamInterceptor", slog.String("full_method", info.FullMethod))
		return handler(srv, ss)
	}
}
