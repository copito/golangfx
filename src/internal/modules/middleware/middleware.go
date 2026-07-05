package middleware

import (
	"go.uber.org/fx"

	"github.com/copito/runner/src/pkg/middleware/logging"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(logging.NewLoggingInterceptor, fx.ResultTags(`group:"grpc_middleware`)),
		fx.Annotate(logging.NewLoggingInterceptor, fx.ResultTags(`group:"grpc_middleware`)),
	),
)
