package logging

import (
	"context"
	"log/slog"

	"github.com/copito/runner/src/internal/entities"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

type LoggingInterceptor struct {
	logger *slog.Logger

	opts       []logging.Option
	logTraceID func(ctx context.Context) logging.Fields
}

func NewLoggingInterceptor(logger *slog.Logger, config *entities.Config) *LoggingInterceptor {
	logTraceID := func(ctx context.Context) logging.Fields {
		span := trace.SpanContextFromContext(ctx)
		if span.IsSampled() {
			return logging.Fields{"traceID", span.TraceID().String()}
		}
		return nil
	}

	// Create the gRPC server object
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other options you want to use here
	}

	rpcLogger := logger.With("service", config.Global.Service)
	return &LoggingInterceptor{logger: rpcLogger, opts: opts, logTraceID: logTraceID}
}

func (l *LoggingInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(InterceptorLogger(l.logger), logging.WithFieldsFromContext(l.logTraceID))
}

func (l *LoggingInterceptor) BuildStreamInterceptor() grpc.StreamServerInterceptor {
	return logging.StreamServerInterceptor(InterceptorLogger(l.logger), logging.WithFieldsFromContext(l.logTraceID))
}
