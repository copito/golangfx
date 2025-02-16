package auth

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
	noToken            = status.Errorf(codes.Unauthenticated, "authorization token is not supplied")
)

type AuthInterceptor struct {
	logger       *slog.Logger
	authProvider AuthProvider
}

func NewAuthInterceptor(env string, logger *slog.Logger) *AuthInterceptor {
	if env == "local" {
		provider := NewLocalAuthProvider(logger)
		return &AuthInterceptor{authProvider: provider, logger: logger}
	} else {
		provider := NewBasicAuthProvider(logger, nil)
		return &AuthInterceptor{authProvider: provider, logger: logger}
	}
}

// valid validates the authorization
func (a AuthInterceptor) CaptureValidUser(authorization []string) (string, bool) {
	if len(authorization) < 1 {
		return "", false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string
	isValid := a.authProvider.IsValid(token)
	if !isValid {
		return "", false
	}

	user, err := a.authProvider.GetUser(token)
	if err != nil {
		return "", false
	}

	return user, true
}

func (a AuthInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		a.logger.Debug("middleware: authentication checker", slog.String("full_method", info.FullMethod))

		// authentication (token verification)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}

		authHeader, ok := md["authorization"]
		if !ok {
			return nil, noToken
		}

		user, isValid := a.CaptureValidUser(authHeader)
		if !isValid {
			return nil, errInvalidToken
		}

		// Add user to context
		// ctx = middleware.AppendToInterceptContext(ctx, "user", []string{user})
		ctx = context.WithValue(ctx, "username", user) //go:lint SA1029

		// Run function
		m, err := handler(ctx, req)
		if err != nil {
			a.logger.Error("RPC failed with error", slog.Any("err", err))
		}
		return m, err
	}
}

func (a AuthInterceptor) BuildStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// authentication (token verification)
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return errMissingMetadata
		}

		authHeader, ok := md["authorization"]
		if !ok {
			return noToken
		}

		user, isValid := a.CaptureValidUser(authHeader)
		if !isValid {
			return errInvalidToken
		}

		// Add user to context
		ss.SetHeader(metadata.Pairs("user", user))

		// Run function
		err := handler(srv, ss)
		if err != nil {
			a.logger.Error("RPC failed with error", slog.Any("err", err))
		}
		return err
	}
}
