package handler

import "go.uber.org/fx"

// Module provides all handlers as a group for Fx.
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewHealthCheckHandler, fx.ResultTags(`group:"grpc_handlers"`)),
		fx.Annotate(NewRunnerHandler, fx.ResultTags(`group:"grpc_handlers"`)),
		// Add other handlers here using fx.Annotate.
		// i.e. fx.Annotate(NewBuilderHandler, fx.ResultTags(`group:"grpc_handlers"`)),
		// i.e. fx.Annotate(NewGlobalHandler, fx.ResultTags(`group:"grpc_handlers"`)),
	),
	fx.Provide(
		fx.Annotate(NewMainHandler, fx.As(new(HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
		fx.Annotate(NewSwaggerFileHandler, fx.As(new(HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
		fx.Annotate(NewSwaggerHandler, fx.As(new(HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
		fx.Annotate(NewMetricsHandler, fx.As(new(HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
	),
)
