package ratelimiter

import "go.uber.org/fx"

var Module = fx.Provide(
	NewUserLimiter,
)
