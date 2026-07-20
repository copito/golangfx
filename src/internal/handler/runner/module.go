package runner

import "go.uber.org/fx"

var Module = fx.Provide(
	fx.Annotate(NewRunnerHandler, fx.ResultTags(`group:"grpc_handlers"`)),
)
