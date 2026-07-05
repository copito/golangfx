package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/internal/modules/config"
)

// DBParams defines the dependencies needed by NewGormDB.
type Params struct {
	fx.In

	Lifecycle      fx.Lifecycle
	Logger         *slog.Logger
	ConfigProvider config.ConfigProvider
}

type Result struct {
	fx.Out

	DB *pgxpool.Pool
}

func ConfigDatabase(logger *slog.Logger, config *entities.Config) *pgxpool.Config {
	dbConfigEntity := config.Database

	const defaultMaxConns = int32(10)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour // time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(dbConfigEntity.ConnectionString)
	if err != nil {
		logger.Error(
			"Failed to create a config, error",
			slog.Any("err", err),
		)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	// dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
	// 	logger.Info("Before acquiring the connection pool to the database...")
	// 	return true
	// }

	// dbConfig.AfterRelease = func(c *pgx.Conn) bool {
	// 	logger.Info("After releasing the connection pool to the database...")
	// 	return true
	// }

	// dbConfig.BeforeClose = func(c *pgx.Conn) {
	// 	logger.Info("Closed the connection pool to the database...")
	// }

	return dbConfig
}

// NewDatabase initializes a PGXPool connection to Postgres with lifecycle management.
func NewDatabasePool(params Params) (Result, error) {
	params.Logger.Info("setting up Database module (with PGXPool)...")

	ctx := context.Background()

	// conn, err := pgxpool.New(ctx, dbConfig.ConnectionString)
	// if err != nil {
	// 	return DBPoolResults{}, err
	// }

	config := params.ConfigProvider.Get()

	dbConfig := ConfigDatabase(params.Logger, config)
	conn, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return Result{}, err
	}

	// conn, err := pgx.Connect(ctx, dbConfig.ConnectionString)
	// if err != nil {
	// 	return DBPoolResults{}, err
	// }

	// connection, err := conn.Acquire(context.Background())
	// if err != nil {
	// 	params.Logger.Info("Error while acquiring connection from the database pool!!")
	// }
	// defer connection.Release()

	// Use fx lifecycle hooks to manage the database connection
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.Logger.Info("Database connection established")
			err := conn.Ping(ctx)
			return err
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Closing database connection")
			conn.Close()
			return nil
		},
	})

	return Result{DB: conn}, nil
}
