package health

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(NewHealthHandler, fx.ResultTags(`group:"grpc_handlers"`)),
)
