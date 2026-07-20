package gateway

import (
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/gateway/chucknorris"
)

var Module = fx.Options(
	chucknorris.Module,
)
