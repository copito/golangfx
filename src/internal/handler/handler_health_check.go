package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthCheckHandler struct {
	// Required by GRPC
	health.UnimplementedHealthServer

	// Common parameters that will be used by Handlers
	CommonHandlerParams
}

// Has to always return the interface GRPCHandlerInterface in the signature
// to be picked up by the fx framework and added to the list of handlers.
// Please always add to handlers.go in the modules packages too
func NewHealthCheckHandler(params CommonHandlerParams) GRPCHandlerInterface {
	return &HealthCheckHandler{CommonHandlerParams: params}
}

// Register registers the RunnerHandlers to the gRPC server.
func (s *HealthCheckHandler) RegisterGRPC(server *grpc.Server) {
	s.Logger.Info("registering GRPC handler", slog.String("handler", "health_check"), slog.String("type", "grpc"))
	health.RegisterHealthServer(server, s)
}

func (s *HealthCheckHandler) RegisterGRPCGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	s.Logger.Info("registering GRPC-Gateway handler", slog.String("handler", "health_check"), slog.String("type", "grpc"))
	healthClient := health.NewHealthClient(conn)

	// Register the health check handler
	// https:// grpc-ecosystem.github.io/grpc-gateway/docs/operations/health_check/#adding-healthz-endpoint-to-runtimeservermux
	// Also add to openapi spec - given not proto generated
	mux.HandlePath("GET", "/healthz", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Perform the gRPC health check
		response, err := healthClient.Check(ctx, &health.HealthCheckRequest{})
		if err != nil {
			http.Error(w, "Health check failed", http.StatusInternalServerError)
		}

		// Map gRPC status to an HTTP response
		httpStatus := http.StatusOK
		if response.Status != health.HealthCheckResponse_SERVING {
			httpStatus = http.StatusServiceUnavailable
		}

		// Wrtie the response as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(health.HealthCheckResponse{
			Status: response.Status,
		})
	})

	return nil
}

func (server *HealthCheckHandler) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	// Perform the health check logic here
	// For example, you can check if all dependencies are up and running
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}
