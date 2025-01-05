package modules

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/internal/handler"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCGatewayParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
	Config    *entities.Config
	Handlers  []handler.GRPCHandlerInterface `group:"grpc_handlers"` // Collect all handlers from the group.

	// adding as requirement to force order (dependency)
	GRPCServer *grpc.Server
}

type GRPCGatewayResults struct {
	fx.Out

	Mux        *runtime.ServeMux
	HttpServer *http.Server
}

func NewGRPCGateway(params GRPCGatewayParams) (GRPCGatewayResults, error) {
	params.Logger.Info("setting up gRPC Gateway module...")

	backendConfig := params.Config.Backend

	grpcFullUrl := "0.0.0.0" + backendConfig.GrpcPort
	conn, err := grpc.NewClient(
		grpcFullUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		params.Logger.Error(
			"failed to dial grpc server to setup gRPC-Gateway",
			slog.String("full_url", grpcFullUrl),
			slog.Any("err", err),
		)
		return GRPCGatewayResults{}, fmt.Errorf("failed to open grpc connection to own server for gRPC-Gateway: %w", err)
	}

	mux := runtime.NewServeMux()

	// Dynamically register all handlers
	for _, handler := range params.Handlers {
		if err := handler.RegisterGRPCGateway(context.Background(), mux, conn); err != nil {
			params.Logger.Error(
				"failed to register gRPC-Gateway handler",
				slog.Any("err", err),
			)
			return GRPCGatewayResults{}, fmt.Errorf("failed to register gRPC service handler: %w", err)
		}
	}

	server := &http.Server{
		Addr:    backendConfig.HttpPort,
		Handler: mux,
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				params.Logger.Info("Serving gRPC-Gateway (HTTP) on: " + backendConfig.HttpPort)
				if err := server.ListenAndServe(); err != nil {
					log.Fatalf("Failed to serve HTTP: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Stopping gRPC-Gateway (HTTP) server...")
			return server.Close()
		},
	})

	return GRPCGatewayResults{
		Mux:        mux,
		HttpServer: server,
	}, nil
}

var GRPCGatewayModule = fx.Provide(NewGRPCGateway)
