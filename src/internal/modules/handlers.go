package modules

import (
	"github.com/copito/runner/src/internal/handler"
	"go.uber.org/fx"
)

// Module provides all handlers as a group for Fx.
var HandlerModule = fx.Options(
	fx.Provide(
		fx.Annotate(handler.NewHealthCheckHandler, fx.ResultTags(`group:"grpc_handlers"`)),
		fx.Annotate(handler.NewRunnerHandler, fx.ResultTags(`group:"grpc_handlers"`)),
		// Add other handlers here using fx.Annotate.
		// i.e. fx.Annotate(handler.NewBuilderHandler, fx.ResultTags(`group:"grpc_handlers"`)),
		// i.e. fx.Annotate(handler.NewGlobalHandler, fx.ResultTags(`group:"grpc_handlers"`)),
	),
)

// Module provides all handlers as a group for Fx.
var AdditionalHandlerModule = fx.Options(
	fx.Provide(
		fx.Annotate(handler.NewSwaggerFileHandler, fx.As(new(handler.HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
		fx.Annotate(handler.NewSwaggerHandler, fx.As(new(handler.HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
	),
)
