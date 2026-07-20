package chucknorris

import "go.uber.org/fx"

var Module = fx.Provide(
	NewChuckNorrisGateway,
)