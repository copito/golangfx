package swagger

import (
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/handler/common"
)

var Module = fx.Provide(
	fx.Annotate(NewSwaggerFileHandler, fx.As(new(common.HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
	fx.Annotate(NewSwaggerHandler, fx.As(new(common.HttpHandlerInterface)), fx.ResultTags(`group:"additional_handlers"`)),
)
