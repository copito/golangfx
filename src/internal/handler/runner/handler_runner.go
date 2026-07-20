package runner

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/copito/runner/idl_gen/runner/v1"
	"github.com/copito/runner/src/internal/handler/common"
	"github.com/copito/runner/src/pkg/logger"
)

type RunnerHandler interface {
	common.GRPCHandlerInterface

	Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error)
	GetAvailableEngines(context.Context, *pb.GetAvailableEnginesRequest) (*pb.GetAvailableEnginesResponse, error)
	SubmitQuery(context.Context, *pb.SubmitQueryRequest) (*pb.QueryResultResponse, error)
	SubmitQueryAsync(context.Context, *pb.SubmitQueryRequest) (*pb.SubmitQueryAsyncResponse, error)
	CheckQueryStatus(context.Context, *pb.CheckQueryStatusRequest) (*pb.CheckQueryStatusResponse, error)
	GetQueryResult(*pb.GetQueryResultRequest, grpc.ServerStreamingServer[pb.QueryResultResponse]) error
}

var _ RunnerHandler = (*runnerHandler)(nil)

type Params struct {
	fx.In
	Logger *slog.Logger
}

type Result struct {
	RunnerHandler RunnerHandler
}

type runnerHandler struct {
	// Required by GRPC
	pb.UnimplementedRunnerServiceServer

	// Common parameters that will be used by Handlers
	Logger *slog.Logger
}

// Has to always return the interface GRPCHandlerInterface in the signature
// to be picked up by the fx framework and added to the list of handlers.
// Please always add to handlers.go in the modules packages too
func NewRunnerHandler(params Params) Result {
	return Result{
		RunnerHandler: &runnerHandler{
			Logger: params.Logger,
		},
	}
}

// Register registers the RunnerHandlers to the gRPC server.
func (s *runnerHandler) RegisterGRPC(server *grpc.Server) {
	s.Logger.Info("registering GRPC handler", slog.String("handler", "runner"), slog.String("type", "grpc"))
	pb.RegisterRunnerServiceServer(server, s)
}

func (s *runnerHandler) RegisterGRPCGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	s.Logger.Info("registering GRPC-Gateway handler", slog.String("handler", "runner"), slog.String("type", "grpc"))
	return pb.RegisterRunnerServiceHandler(ctx, mux, conn)
}

func (s runnerHandler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	logger := logger.LoggerFromContext(ctx, slog.Default())
	logger.Info("responsing to ping test...", slog.String("path", "ping"))
	return &pb.PingResponse{
		Message: "PONG",
	}, nil
}

func (s runnerHandler) GetAvailableEngines(context.Context, *pb.GetAvailableEnginesRequest) (*pb.GetAvailableEnginesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvailableEngines not implemented")
}

func (s runnerHandler) SubmitQuery(context.Context, *pb.SubmitQueryRequest) (*pb.QueryResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitQuery not implemented")
}

func (s runnerHandler) SubmitQueryAsync(context.Context, *pb.SubmitQueryRequest) (*pb.SubmitQueryAsyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitQueryAsync not implemented")
}

func (s runnerHandler) CheckQueryStatus(context.Context, *pb.CheckQueryStatusRequest) (*pb.CheckQueryStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckQueryStatus not implemented")
}

func (s runnerHandler) GetQueryResult(*pb.GetQueryResultRequest, grpc.ServerStreamingServer[pb.QueryResultResponse]) error {
	return status.Errorf(codes.Unimplemented, "method GetQueryResult not implemented")
}
