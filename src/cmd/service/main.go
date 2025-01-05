package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/copito/runner/src/internal/modules"
	"github.com/copito/runner/src/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Format fx logger
		// fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
		// 	return &fxevent.SlogLogger{Logger: log}
		// }),
		modules.LoggerModule,
		modules.ConfigModule,
		modules.DatabaseModule,
		modules.DatabasePoolModule,
		modules.RepositoryModule,
		modules.KafkaProducerModule,
		modules.KafkaConsumerModule,
		modules.HandlerModule,
		modules.GrpcServerModule,
		modules.GRPCGatewayModule,

		// fx.Invoke(func(logger *slog.Logger, config *entities.Config, db *gorm.DB) {
		// 	logger.Info(
		// 		"Testing",
		// 		slog.Bool("is", true),
		// 		slog.Time("time", time.Now()),
		// 		slog.Any("example", config.Database.ConnectionString),
		// 	)

		// 	type Result struct {
		// 		Value int
		// 	}

		// 	var result Result
		// 	tx := db.Raw("SELECT 1 as value;").Scan(&result)
		// 	logger.Info(
		// 		"Result",
		// 		slog.Any("result", result.Value),
		// 		slog.Any("tx", tx.Error),
		// 	)
		// }),
		// fx.Invoke(fake_publisher.GenerateFakeData),
		// fx.Invoke(func(GrpcServer *grpc.Server, grpcGatewayServer *http.Server, repo *repository.Queries) {

		fx.Invoke(func(logger *slog.Logger, repo *repository.Queries, connPool *pgxpool.Pool) {
			fmt.Print("testing")
			ctx := context.Background()

			conn, err := connPool.Acquire(ctx)
			if err != nil {
				logger.Info("Error while acquiring connection from the database pool!!")
			}
			defer conn.Release()

			users, err := repo.GetAllUsers(ctx, conn)
			if err != nil {
				logger.Error("unable to get all users", slog.Any("err", err))
				panic("error!!")
			}

			for _, user := range users {
				logger.Info(fmt.Sprintf("%+v\n", user))
				logger.Info(fmt.Sprintf("Email: %v\n", user.Email))
				logger.Info(fmt.Sprintf("ID: %v\n", user.ID))
				logger.Info(fmt.Sprintf("Username: %v\n", user.Username))
				logger.Info(fmt.Sprintf("Created At: %v\n", user.CreatedAt.Time))
			}

			logger.Info("Successfully fetched all users")
		}),
	).Run()
}
