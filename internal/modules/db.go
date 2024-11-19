package modules

import (
	"context"
	"log/slog"

	"github.com/copito/runner/internal/entities"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBParams defines the dependencies needed by NewGormDB.
type DBParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
	Config    *entities.Config
}

type DBResults struct {
	fx.Out

	DB *gorm.DB
}

// NewDatabase initializes a GORM DB connection to Postgres with lifecycle management.
func NewDatabase(params DBParams) (DBResults, error) {
	params.Logger.Info("setting up Database (with GORM)...")
	dbConfig := params.Config.Database

	// Set up GORM with a customized logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		// Logger: params.Logger,
	}

	db, err := gorm.Open(postgres.Open(dbConfig.ConnectionString), gormConfig)
	if err != nil {
		return DBResults{}, err
	}

	// Set up connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return DBResults{}, err
	}
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	// Use fx lifecycle hooks to manage the database connection
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.Logger.Info("Database connection established")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Closing database connection")
			return sqlDB.Close()
		},
	})

	return DBResults{DB: db}, nil
}

var DatabaseModule = fx.Provide(NewDatabase)
