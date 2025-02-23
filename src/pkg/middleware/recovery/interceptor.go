package recovery

import (
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/copito/runner/src/internal/entities"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
)

type RecoveryInterceptor struct {
	logger *slog.Logger

	opts []recovery.Option
}

func NewRecoveryInterceptor(logger *slog.Logger, ms *entities.MetricStore) *RecoveryInterceptor {
	// Define panic recovery handler (if metric store is provided or not provided)
	var grpcPanicRecoveryHandler func(p any) (err error)
	if ms == nil {
		// Define panic recovery handler
		grpcPanicRecoveryHandler = func(p any) (err error) {
			logger.Error("recovered from panic", slog.Any("panic", p), slog.String("stack", string(debug.Stack())))
			return status.Errorf(codes.Internal, "%s", p)
		}
	} else {
		// Define panic recovery handler with metric store
		grpcPanicRecoveryHandler = func(p any) (err error) {
			ms.PanicsTotal.Inc()
			logger.Error("recovered from panic", slog.Any("panic", p), slog.String("stack", string(debug.Stack())))
			return status.Errorf(codes.Internal, "%s", p)
		}
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(grpcPanicRecoveryHandler),
	}

	return &RecoveryInterceptor{logger: logger, opts: recoveryOpts}
}

func (l RecoveryInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	return recovery.UnaryServerInterceptor(l.opts...)
}

func (l RecoveryInterceptor) BuildStreamInterceptor() grpc.StreamServerInterceptor {
	return recovery.StreamServerInterceptor(l.opts...)
}
