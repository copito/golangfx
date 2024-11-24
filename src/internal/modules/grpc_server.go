package modules

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/copito/runner/src/internal/entities"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/copito/runner/idl_gen/runner/v1"
)

type server struct {
	pb.UnimplementedRunnerServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message: "PONG",
	}, nil
}

func (s server) GetAvailableEngines(context.Context, *pb.GetAvailableEnginesRequest) (*pb.GetAvailableEnginesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAvailableEngines not implemented")
}

func (s server) SubmitQuery(context.Context, *pb.SubmitQueryRequest) (*pb.QueryResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitQuery not implemented")
}

func (s server) SubmitQueryAsync(context.Context, *pb.SubmitQueryRequest) (*pb.SubmitQueryAsyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitQueryAsync not implemented")
}

func (s server) CheckQueryStatus(context.Context, *pb.CheckQueryStatusRequest) (*pb.CheckQueryStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckQueryStatus not implemented")
}

func (s server) GetQueryResult(*pb.GetQueryResultRequest, grpc.ServerStreamingServer[pb.QueryResultResponse]) error {
	return status.Errorf(codes.Unimplemented, "method GetQueryResult not implemented")
}

// DBParams defines the dependencies needed by NewGormDB.
type GRPCParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
	Config    *entities.Config
}

type GRPCResults struct {
	fx.Out

	GrpcServer *grpc.Server
}

// NewGRPCServer initializes a grpc server using proto files generated
func NewGRPCServer(params GRPCParams) (GRPCResults, error) {
	params.Logger.Info("setting up gRPC (with gRPC + Proto + gRPC Gateway)...")
	backendConfig := params.Config.Backend

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", backendConfig.GrpcPort)
	if err != nil {
		params.Logger.Error(
			"Failed to open listener for grpc",
			slog.String("port", backendConfig.GrpcPort),
			slog.Any("err", err),
		)
		panic("failed to listen to grpc port")
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	pb.RegisterRunnerServiceServer(s, &server{})

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
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
		panic("unable to open grpc connection to own server for gRPC-Gateway")
	}

	gwmux := runtime.NewServeMux()
	// Register Runner
	err = pb.RegisterRunnerServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		params.Logger.Error(
			"failed to register gateway gRPC-Gateway",
			slog.String("handler", "RegisterRunnerServiceHandler"),
			slog.Any("err", err),
		)
		panic("unable to register RegisterRunnerServiceHandler handler...")
	}

	gwServer := &http.Server{
		Addr:    backendConfig.HttpPort,
		Handler: gwmux,
	}

	// Use fx lifecycle hooks to manage the database connection
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Serve gRPC server
			go func(lis net.Listener) {
				params.Logger.Info("Serving gRPC on: " + grpcFullUrl)
				log.Fatalln(s.Serve(lis))
			}(lis)

			go func(gwServer *http.Server) {
				params.Logger.Info("Serving gRPC-Gateway on: http://0.0.0.0" + backendConfig.HttpPort)
				log.Fatalln(gwServer.ListenAndServe())
			}(gwServer)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Closing gRPC connection for port " + backendConfig.GrpcPort)
			s.GracefulStop()

			params.Logger.Info("Closing gRPC connection for port " + backendConfig.HttpPort)
			_ = gwServer.Close()
			return nil
		},
	})

	return GRPCResults{GrpcServer: s}, nil
}

var GRPCModule = fx.Provide(NewGRPCServer)
