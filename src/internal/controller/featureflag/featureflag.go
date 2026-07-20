package featureflag

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/modules/config"
)

type FeatureFlagger interface {
	AllowDatabaseType(dbName string) bool
}

var _ FeatureFlagger = (*featureFlagger)(nil)

type Params struct {
	fx.In
	Lifecycle      fx.Lifecycle
	Logger         slog.Logger
	ConfigProvider config.ConfigProvider
}

type Result struct {
	fx.Out
	FeatureFlagger FeatureFlagger
}

type featureFlagger struct {
	ConfigProvider config.ConfigProvider
}

func NewFeatureFlagger(params Params) Result {
	featureFlagger := &featureFlagger{
		ConfigProvider: params.ConfigProvider,
	}

	return Result{
		FeatureFlagger: featureFlagger,
	}
}

func (f *featureFlagger) AllowDatabaseType(dbName string) bool {
	// Some logic here for example
	return true
}
