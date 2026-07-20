package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/copito/runner/src/internal/handler/common"
)

type HealthHandler interface {
	common.GRPCHandlerInterface

	Check(ctx context.Context, req *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error)
}

var _ HealthHandler = (*healthHandler)(nil)

type Params struct {
	fx.In

	Logger *slog.Logger
}

type healthHandler struct {
	// Required by GRPC
	healthv1.UnimplementedHealthServer

	Logger *slog.Logger
}

// Has to always return the interface GRPCHandlerInterface in the signature
// to be picked up by the fx framework and added to the list of handlers.
// Please always add to handlers.go in the modules packages too
func NewHealthHandler(logger *slog.Logger) common.GRPCHandlerInterface {
	return &healthHandler{
		Logger: logger,
	}
}

// Register registers the RunnerHandlers to the gRPC server.
func (s *healthHandler) RegisterGRPC(server *grpc.Server) {
	s.Logger.Info("registering GRPC handler", slog.String("handler", "health_check"), slog.String("type", "grpc"))
	healthv1.RegisterHealthServer(server, s)
}

func (s *healthHandler) RegisterGRPCGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	s.Logger.Info("registering GRPC-Gateway handler", slog.String("handler", "health_check"), slog.String("type", "grpc"))
	healthClient := healthv1.NewHealthClient(conn)

	// Register the health check handler
	// https:// grpc-ecosystem.github.io/grpc-gateway/docs/operations/health_check/#adding-healthz-endpoint-to-runtimeservermux
	// Also add to openapi spec - given not proto generated
	mux.HandlePath("GET", "/healthz", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Perform the gRPC health check
		response, err := healthClient.Check(ctx, &healthv1.HealthCheckRequest{})
		if err != nil {
			http.Error(w, "Health check failed", http.StatusInternalServerError)
		}

		// Map gRPC status to an HTTP response
		httpStatus := http.StatusOK
		if response.Status != healthv1.HealthCheckResponse_SERVING {
			httpStatus = http.StatusServiceUnavailable
		}

		// Write the response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(healthv1.HealthCheckResponse{
			Status: response.Status,
		})
	})

	return nil
}

func (server *healthHandler) Check(ctx context.Context, req *healthv1.HealthCheckRequest) (*healthv1.HealthCheckResponse, error) {
	// Perform the health check logic here
	// For example, you can check if all dependencies are up and running
	return &healthv1.HealthCheckResponse{
		Status: healthv1.HealthCheckResponse_SERVING,
	}, nil
}
