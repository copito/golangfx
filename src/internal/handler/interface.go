package handler

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const validPatternVersion = 1

// ServiceHandler defines an interface for gRPC service handlers.
type GRPCHandlerInterface interface {
	RegisterGRPC(server *grpc.Server)
	RegisterGRPCGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

// Define HTTP handler methods here
type HttpHandlerInterface interface {
	ServeHTTP() runtime.HandlerFunc
	Method() string
	Pattern() runtime.Pattern
}
