package handler

import (
	"context"
	"log/slog"

	pb "github.com/copito/runner/idl_gen/runner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RunnerHandler struct {
	// Required by GRPC
	pb.UnimplementedRunnerServiceServer

	// Common parameters that will be used by Handlers
	CommonHandlerParams
}

// Has to always return the interface GRPCHandlerInterface in the signature
// to be picked up by the fx framework and added to the list of handlers.
// Please always add to handlers.go in the modules packages too
func NewRunnerHandler(params CommonHandlerParams) GRPCHandlerInterface {
	return &RunnerHandler{CommonHandlerParams: params}
}

// Register registers the RunnerHandlers to the gRPC server.
func (s *RunnerHandler) Register(server *grpc.Server) {
	s.Logger.Info("registering handler", slog.String("handler", "runner"), slog.String("type", "grpc"))
	pb.RegisterRunnerServiceServer(server, s)
}

func (s RunnerHandler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	s.Logger.Info("responsing to ping test...", slog.String("path", "ping"))
	return &pb.PingResponse{
		Message: "PONG",
	}, nil
}

func (s RunnerHandler) GetAvailableEngines(context.Context, *pb.GetAvailableEnginesRequest) (*pb.GetAvailableEnginesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvailableEngines not implemented")
}

func (s RunnerHandler) SubmitQuery(context.Context, *pb.SubmitQueryRequest) (*pb.QueryResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitQuery not implemented")
}

func (s RunnerHandler) SubmitQueryAsync(context.Context, *pb.SubmitQueryRequest) (*pb.SubmitQueryAsyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitQueryAsync not implemented")
}

func (s RunnerHandler) CheckQueryStatus(context.Context, *pb.CheckQueryStatusRequest) (*pb.CheckQueryStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckQueryStatus not implemented")
}

func (s RunnerHandler) GetQueryResult(*pb.GetQueryResultRequest, grpc.ServerStreamingServer[pb.QueryResultResponse]) error {
	return status.Errorf(codes.Unimplemented, "method GetQueryResult not implemented")
}
