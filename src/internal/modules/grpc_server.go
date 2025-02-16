package modules

import (
	"context"
	"log"
	"log/slog"
	"net"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/internal/handler"
	"github.com/copito/runner/src/pkg/middleware"
	"github.com/copito/runner/src/pkg/middleware/auth"
	authbypass "github.com/copito/runner/src/pkg/middleware/auth_bypass"
	info "github.com/copito/runner/src/pkg/middleware/info"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Params contains the dependencies for creating the gRPC server.
type GrpcParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
	Config    *entities.Config
	Handlers  []handler.GRPCHandlerInterface `group:"grpc_handlers"` // Collect all handlers from the group.
}

// Results is the output of the gRPC server module.
type GrpcResults struct {
	fx.Out

	GrpcServer *grpc.Server
}

// NewGRPCServer initializes the gRPC server.
func NewGRPCServer(params GrpcParams) (GrpcResults, error) {
	params.Logger.Info("setting up gRPC server module...")

	backendConfig := params.Config.Backend
	listener, err := net.Listen("tcp", backendConfig.GrpcPort)
	if err != nil {
		params.Logger.Error(
			"Failed to open listener for gRPC",
			slog.String("port", backendConfig.GrpcPort),
			slog.Any("err", err),
		)
		return GrpcResults{}, err
	}

	// Create the gRPC server object.
	authBypassInterceptor := authbypass.NewLocalBypassAuthInterceptor(params.Logger, params.Config)
	authInterceptor := auth.NewAuthInterceptor(params.Config.Backend.Environment, params.Logger)
	infoInterceptor := info.NewInfoInterceptor(params.Logger)

	unaryInterceptors := grpc.UnaryInterceptor(
		middleware.ChainUnaryInterceptors(
			infoInterceptor.BuildUnaryInterceptor(),
			authBypassInterceptor.BuildUnaryInterceptor(),
			authInterceptor.BuildUnaryInterceptor(),
		),
	)

	// Create the gRPC server object.
	server := grpc.NewServer(unaryInterceptors)

	// Register all service handlers.
	for _, handler := range params.Handlers {
		handler.RegisterGRPC(server)
	}

	// Start the server using Fx lifecycle hooks.
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			params.Logger.Info("serving gRPC on port " + backendConfig.GrpcPort)
			go func() {
				if err := server.Serve(listener); err != nil {
					params.Logger.Error("failed to serve gRPC", slog.Any("err", err))
					log.Fatalf("Failed to serve gRPC: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Stopping gRPC server...")
			server.GracefulStop()
			return nil
		},
	})

	return GrpcResults{GrpcServer: server}, nil
}

var GrpcServerModule = fx.Provide(NewGRPCServer)
