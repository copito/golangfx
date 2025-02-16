package log

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

type InfoInterceptor struct {
	logger *slog.Logger
}

func NewInfoInterceptor(logger *slog.Logger) *InfoInterceptor {
	return &InfoInterceptor{logger: logger}
}

// UnaryInterceptor is a unary Interceptor that logs the method name and request
func (i InfoInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		i.logger.Info("received UnaryInterceptor", slog.String("full_method", info.FullMethod), slog.Any("req", req))
		resp, err := handler(ctx, req)
		i.logger.Info("ended UnaryInterceptor", slog.String("full_method", info.FullMethod), slog.Any("res", resp))
		return resp, err
	}
}

// StreamInterceptor is a stream interceptor that logs the method name
func (i InfoInterceptor) BuildStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		i.logger.Info("running StreamInterceptor", slog.String("full_method", info.FullMethod))

		err := handler(srv, ss)

		if err != nil {
			i.logger.Error("Stream error", slog.String("full_method", info.FullMethod), slog.Any("err", err))
		} else {
			i.logger.Info("Stream finished successfully", slog.String("full_method", info.FullMethod))
		}

		return err
	}
}
