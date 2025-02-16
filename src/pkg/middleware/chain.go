package middleware

import (
	"context"

	"google.golang.org/grpc"
)

// ChainUnaryInterceptors chains multiple interceptors into a single interceptor.
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// Start with the last interceptor in the chain, which calls the handler.
		chainedHandler := handler

		// Chain the interceptors in reverse order. (so the first is the last to be called)
		for i := len(interceptors) - 1; i >= 0; i-- {
			currentInterceptor := interceptors[i]
			nextHandler := chainedHandler

			chainedHandler = func(ctx context.Context, currentReq any) (any, error) {
				return currentInterceptor(ctx, currentReq, info, nextHandler)
			}
		}

		return chainedHandler(ctx, req)
	}
}

// ChainStreamInterceptors chains multiple stream interceptors into a single interceptor
func ChainStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Start with the last interceptor in the chain, which calls the handler.
		chainedHandler := handler

		// Chain the interceptors in reverse order. (so the first is the last to be called)
		for i := len(interceptors) - 1; i >= 0; i-- {
			currentInterceptor := interceptors[i]
			nextHandler := chainedHandler

			chainedHandler = func(currentSrv any, currentStream grpc.ServerStream) error {
				return currentInterceptor(currentSrv, currentStream, info, nextHandler)
			}
		}

		return chainedHandler(srv, ss)
	}
}
