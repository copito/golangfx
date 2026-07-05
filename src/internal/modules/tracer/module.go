package tracer

import "go.uber.org/fx"

var Module = fx.Provide(NewTraceProvider)
