package controller

import (
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/controller/featureflag"
	"github.com/copito/runner/src/internal/controller/ratelimiter"
)

var Module = fx.Options(
	// Support Modules
	ratelimiter.Module,
	featureflag.Module,

	// fake_publisher.Module,
)
