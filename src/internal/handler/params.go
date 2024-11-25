package handler

import (
	"log/slog"

	"github.com/copito/runner/src/internal/entities"
	"go.uber.org/fx"
)

type CommonHandlerParams struct {
	fx.In

	Logger *slog.Logger
	Config *entities.Config
}
