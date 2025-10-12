package validate

import (
	"context"
	"log/slog"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type ProtoValidatorInterceptor struct {
	logger *slog.Logger
}

func NewProtoValidatorInterceptor(logger *slog.Logger) ProtoValidatorInterceptor {
	return ProtoValidatorInterceptor{logger: logger}
}

func (i ProtoValidatorInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// proto validate the request
		message, ok := req.(proto.Message)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "expected proto.Message, got %T", req)
		}
		err := protovalidate.Validate(message)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}

		resp, err := handler(ctx, req)
		return resp, err
	}
}

func (i ProtoValidatorInterceptor) BuildStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// proto validate the request
		message, ok := srv.(proto.Message)
		if !ok {
			return status.Errorf(codes.InvalidArgument, "expected proto.Message, got %T", srv)
		}
		err := protovalidate.Validate(message)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}

		err = handler(srv, ss)
		return err
	}
}
