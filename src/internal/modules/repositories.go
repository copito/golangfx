package modules

import (
	"context"
	"log/slog"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type RepositoryParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
	Config    *entities.Config
	DB        *pgxpool.Pool
}

type RepositoryResults struct {
	fx.Out

	Repositories *repository.Queries
}

func NewRepository(params RepositoryParams) (RepositoryResults, error) {
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

	return RepositoryResults{Repositories: repo}, nil
}

var RepositoryModule = fx.Provide(NewRepository)
