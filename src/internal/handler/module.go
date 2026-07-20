package handler

import (
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/handler/health"
	"github.com/copito/runner/src/internal/handler/mainhandler"
	"github.com/copito/runner/src/internal/handler/metrics"
	"github.com/copito/runner/src/internal/handler/runner"
	"github.com/copito/runner/src/internal/handler/swagger"
)

// Module provides all handlers as a group for Fx.
var Module = fx.Options(
	health.Module,
	mainhandler.Module,
	metrics.Module,
	swagger.Module,
	runner.Module,
)
