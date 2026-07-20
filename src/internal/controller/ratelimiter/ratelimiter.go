package ratelimiter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/modules/config"
)

type UserLimiter interface {
	Allow(userID string) bool
	AllowN(UserID string, n int) bool
}

var _ UserLimiter = (*userLimiter)(nil)

type Params struct {
	fx.In

	Lifecycle      fx.Lifecycle
	Logger         *slog.Logger
	ConfigProvider config.ConfigProvider
}

type Result struct {
	fx.Out

	UserLimiter UserLimiter
}

type userLimiter struct {
	user map[string]*TokenBucket
	mu   sync.RWMutex

	bucketCapacity       int
	refillPerSecond      float64
	cleanupIntervalInSec float64

	stopChan chan struct{}
}

func NewUserLimiter(params Params) Result {
	// config := params.ConfigProvider.Get()

	// normally coming from config
	rateCapacity := 1
	refillPerSecond := 2
	cleanupIntervalInSec := 3
	cleanupMaximumAgeInSec := 4

	limiter := &userLimiter{
		user:                 make(map[string]*TokenBucket),
		bucketCapacity:       rateCapacity,
		refillPerSecond:      float64(refillPerSecond),
		cleanupIntervalInSec: float64(cleanupIntervalInSec),
		stopChan:             make(chan struct{}),
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.Logger.Info("starting user limiter...", slog.Int("capacity", limiter.bucketCapacity))

			go func() {
				params.Logger.Info("starting user limiter cleanup routine....")
				limiter.Cleanup(time.Duration(cleanupMaximumAgeInSec) * time.Second)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("stopping user limiter cleanup routine...")
			limiter.stopChan <- struct{}{}
			return nil
		},
	})

	return Result{
		UserLimiter: limiter,
	}
}

func (cl *userLimiter) GetLimiter(userID string) *TokenBucket {
	cl.mu.RLock()
	limiter, exists := cl.user[userID]
	cl.mu.RUnlock()

	if exists {
		return limiter
	}

	cl.mu.Lock()
	defer cl.mu.Unlock()

	// Double check after acquiring write lock
	if limiter, exists = cl.user[userID]; exists {
		return limiter
	}

	limiter = NewTokenBucket(cl.bucketCapacity, cl.refillPerSecond)
	cl.user[userID] = limiter
	return limiter
}

func (cl *userLimiter) Allow(userID string) bool {
	l := cl.GetLimiter(userID)
	if l == nil {
		return true
	}
	return l.Allow()
}

func (cl *userLimiter) AllowN(userID string, n int) bool {
	l := cl.GetLimiter(userID)
	if l == nil {
		return true
	}
	return l.AllowN(n)
}

func (cl *userLimiter) Cleanup(maxAge time.Duration) {
	// Remove old clients periodically
	// Track last access time and remove inactive users to prevent unbounded memory growth
	timer := time.NewTicker(time.Duration(cl.cleanupIntervalInSec) * time.Second)
	for {
		select {
		case <-timer.C:
			cl.mu.Lock()
			for userId, limiter := range cl.user {
				if time.Since(limiter.lastRefill) > maxAge {
					delete(cl.user, userId)
				}
			}
			cl.mu.Unlock()
		case <-cl.stopChan:
			timer.Stop()
			return
		}
	}
}
