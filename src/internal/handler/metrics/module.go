package metrics

import (
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/handler/common"
)

var Module = fx.Provide(
	fx.Annotate(NewMetricsHandler, fx.As(new(common.HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
)
