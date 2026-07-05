package main

import (
	"net/http"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/copito/runner/src/internal/handler"
	"github.com/copito/runner/src/internal/modules"
	"github.com/copito/runner/src/internal/modules/config"
	"github.com/copito/runner/src/internal/modules/db"
	"github.com/copito/runner/src/internal/modules/log"
	"github.com/copito/runner/src/internal/modules/metricstore"
	"github.com/copito/runner/src/internal/modules/repo"
	"github.com/copito/runner/src/internal/modules/server"
	"github.com/copito/runner/src/internal/modules/tracer"
)

func main() {
	fx.New(
		// Format fx logger
		// fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
		// 	return &fxevent.SlogLogger{Logger: log}
		// }),
		log.Module,
		config.Module,
		db.Module,
		repo.Module,
		metricstore.Module,
		tracer.Module,
		modules.KafkaProducerModule,
		modules.KafkaConsumerModule,
		handler.Module,
		server.Module,

		fx.Invoke(func(GrpcServer *grpc.Server, grpcGatewayServer *http.Server) {}),
	).Run()
}
