package modules

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/entities"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
}

type MetricResults struct {
	fx.Out

	MetricRegistry *entities.MetricStore
}

func NewMetricStore(params MetricParams) (MetricResults, error) {
	params.Logger.Info("setting up Metric Store module...")

	reg := prometheus.NewRegistry()
	ms := entities.MetricStore{
		Registry: reg,
		PanicsTotal: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "grpc_req_panics_recovered_total",
			Help: "Total number of gRPC requests recovered from internal panic.",
		}),
	}

	return MetricResults{MetricRegistry: &ms}, nil
}

var MetricStoreModule = fx.Provide(NewMetricStore)
