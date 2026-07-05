package repo

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/repository"
)

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
	DB        *pgxpool.Pool
}

type Result struct {
	fx.Out

	Repositories *repository.Queries
}

func NewRepository(params Params) (Result, error) {
	params.Logger.Info("setting up Repository module...")

	// repo := repository.New(params.DB)
	repo := repository.New()

	// Use fx lifecycle hooks to manage the database connection
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.Logger.Info("Setting up repository module...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Closing repository module...")
			return nil
		},
	})

	return Result{Repositories: repo}, nil
}
