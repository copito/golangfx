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
	"github.com/copito/runner/src/pkg/middleware/limiter"
	"github.com/copito/runner/src/pkg/middleware/logging"
	"github.com/copito/runner/src/pkg/middleware/metrics"
	"github.com/copito/runner/src/pkg/middleware/recovery"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// Params contains the dependencies for creating the gRPC server.
type GrpcParams struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Logger      *slog.Logger
	Config      *entities.Config
	MetricStore *entities.MetricStore
	Handlers    []handler.GRPCHandlerInterface `group:"grpc_handlers"` // Collect all handlers from the group.
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
	authInterceptor := auth.NewAuthInterceptor(params.Logger, params.Config)
	infoInterceptor := info.NewInfoInterceptor(params.Logger)
	loggingInterceptor := logging.NewLoggingInterceptor(params.Logger, params.Config)
	limiterInterceptor := limiter.NewLimiterInterceptor(params.Logger)
	recoveryInterceptor := recovery.NewRecoveryInterceptor(params.Logger, params.MetricStore)
	metricsInterceptor := metrics.NewMetricInterceptor(params.Logger, params.Config, params.MetricStore)

	unaryInterceptors := grpc.UnaryInterceptor(
		middleware.ChainUnaryInterceptors(
			recoveryInterceptor.BuildUnaryInterceptor(),
			loggingInterceptor.BuildUnaryInterceptor(),
			infoInterceptor.BuildUnaryInterceptor(),
			limiterInterceptor.BuildUnaryInterceptor(),
			metricsInterceptor.BuildUnaryInterceptor(),
			authBypassInterceptor.BuildUnaryInterceptor(),
			authInterceptor.BuildUnaryInterceptor(),
		),
	)

	// streamInterceptor := grpc.StreamInterceptor(
	// 	middleware.ChainStreamInterceptors(
	// 		recoveryInterceptor.BuildStreamInterceptor(),
	// 		loggingInterceptor.BuildStreamInterceptor(),
	// 		infoInterceptor.BuildStreamInterceptor(),
	// 		limiterInterceptor.BuildStreamInterceptor(),
	// 		metricsInterceptor.BuildStreamInterceptor(),
	// 		authBypassInterceptor.BuildStreamInterceptor(),
	// 		authInterceptor.BuildStreamInterceptor(),
	// 	),
	// )

	// Create the gRPC server object.
	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		unaryInterceptors,
		// streamInterceptor,
	)

	// Register all service handlers.
	for _, handler := range params.Handlers {
		handler.RegisterGRPC(server)
	}

	// Expose to allow users to use grpccurl/postman/etc to instrospect the server
	reflection.Register(server)

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
