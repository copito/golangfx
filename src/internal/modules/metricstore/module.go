package metricstore

import "go.uber.org/fx"

var Module = fx.Provide(NewMetricStore)
