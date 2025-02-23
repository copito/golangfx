package middleware

import "google.golang.org/grpc"

type UnaryInterceptor interface {
	BuildUnaryInterceptor() grpc.UnaryServerInterceptor
}

type StreamInterceptor interface {
	BuildStreamInterceptor() grpc.StreamServerInterceptor
}

type Interceptor interface {
	UnaryInterceptor
	StreamInterceptor
}
