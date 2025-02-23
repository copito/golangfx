package limiter

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/ratelimit"
)

type alwaysPassLimiter struct{}

func (*alwaysPassLimiter) Limit(_ context.Context) error {
	// Example rate limiter could be implemented using e.g. github.com/juju/ratelimit
	//	// Take one token per request. This call doesn't block.
	//	tokenRes := l.tokenBucket.TakeAvailable(1)
	//
	//	// When rate limit reached, return specific error for the clients.
	//	if tokenRes == 0 {
	//		return fmt.Errorf("APP-XXX: reached Rate-Limiting %d", l.tokenBucket.Available())
	//	}
	//
	//	// Rate limit isn't reached.
	//	return nil
	// }
	return nil
}

type LimiterInterceptor struct {
	logger *slog.Logger
}

func NewLimiterInterceptor(logger *slog.Logger) *LimiterInterceptor {
	return &LimiterInterceptor{
		logger: logger,
	}
}

func (l *LimiterInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	limiter := &alwaysPassLimiter{}
	return ratelimit.UnaryServerInterceptor(limiter)
}

func (l *LimiterInterceptor) BuildStreamInterceptor() grpc.StreamServerInterceptor {
	limiter := &alwaysPassLimiter{}
	return ratelimit.StreamServerInterceptor(limiter)
}
