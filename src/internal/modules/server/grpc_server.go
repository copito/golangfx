package server

import (
	"context"
	"log"
	"log/slog"
	"net"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/internal/handler/common"
	"github.com/copito/runner/src/internal/modules/config"
	"github.com/copito/runner/src/pkg/middleware"
	"github.com/copito/runner/src/pkg/middleware/auth"
	authbypass "github.com/copito/runner/src/pkg/middleware/auth_bypass"
	info "github.com/copito/runner/src/pkg/middleware/info"
	"github.com/copito/runner/src/pkg/middleware/limiter"
	"github.com/copito/runner/src/pkg/middleware/logging"

	// "github.com/copito/runner/src/pkg/middleware/metrics"
	"github.com/copito/runner/src/pkg/middleware/recovery"
	"github.com/copito/runner/src/pkg/middleware/tracer"
	"github.com/copito/runner/src/pkg/middleware/validate"
)

// Params contains the dependencies for creating the gRPC server.
type GrpcParams struct {
	fx.In

	Lifecycle      fx.Lifecycle
	Logger         *slog.Logger
	ConfigProvider config.ConfigProvider
	MetricStore    *entities.MetricStore
	MetricServer   *grpcprom.ServerMetrics
	TraceProvider  *sdktrace.TracerProvider
	Handlers       []common.GRPCHandlerInterface `group:"grpc_handlers"` // Collect all handlers from the group.
}

// Results is the output of the gRPC server module.
type GrpcResults struct {
	fx.Out

	GrpcServer *grpc.Server
}

// NewGRPCServer initializes the gRPC server.
func NewGRPCServer(params GrpcParams) (GrpcResults, error) {
	params.Logger.Info("setting up gRPC server module...")
	config := params.ConfigProvider.Get()

	backendConfig := config.Backend
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
	authBypassInterceptor := authbypass.NewLocalBypassAuthInterceptor(params.Logger, config)
	authInterceptor := auth.NewAuthInterceptor(params.Logger, config)
	infoInterceptor := info.NewInfoInterceptor(params.Logger)
	loggingInterceptor := logging.NewLoggingInterceptor(params.Logger, config)
	limiterInterceptor := limiter.NewLimiterInterceptor(params.Logger)
	recoveryInterceptor := recovery.NewRecoveryInterceptor(params.Logger, params.MetricStore)
	metricsInterceptor := tracer.NewTracerInterceptor(params.Logger, config, params.MetricServer) // metrics.NewMetricInterceptor(params.Logger, params.Config, params.MetricStore)
	validatorInterceptor := validate.NewProtoValidatorInterceptor(params.Logger)

	unaryInterceptors := grpc.UnaryInterceptor(
		middleware.ChainUnaryInterceptors(
			recoveryInterceptor.BuildUnaryInterceptor(),
			selector.UnaryServerInterceptor(
				loggingInterceptor.BuildUnaryInterceptor(),
				selector.MatchFunc(loggingInterceptor.ExecuteInterceptor),
			),
			infoInterceptor.BuildUnaryInterceptor(),
			limiterInterceptor.BuildUnaryInterceptor(),
			metricsInterceptor.BuildUnaryInterceptor(),
			authBypassInterceptor.BuildUnaryInterceptor(),
			authInterceptor.BuildUnaryInterceptor(),
			validatorInterceptor.BuildUnaryInterceptor(),
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

			// REQUIRED TO INITIALIZE METRICS
			params.MetricServer.InitializeMetrics(server)

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
