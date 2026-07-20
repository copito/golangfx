package featureflag

import "go.uber.org/fx"

var Module = fx.Provide(
	NewFeatureFlagger,
)