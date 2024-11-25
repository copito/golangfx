package handler

import "google.golang.org/grpc"

// ServiceHandler defines an interface for gRPC service handlers.
type GRPCHandlerInterface interface {
	Register(server *grpc.Server)
}
