package main

import (
	"fmt"
	"net/http"

	"github.com/copito/runner/src/internal/modules"
	"go.uber.org/fx"
	"google.golang.org/grpc"
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
		fx.Invoke(func(GrpcServer *grpc.Server, grpcGatewayServer *http.Server) {
			fmt.Print("testing")
		}),
	).Run()
}
