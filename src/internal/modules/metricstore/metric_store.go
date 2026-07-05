package metricstore

import (
	"log/slog"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/entities"
)

type MetricParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *slog.Logger
}

type MetricResults struct {
	fx.Out

	MetricRegistry *entities.MetricStore
	MetricServer   *grpcprom.ServerMetrics
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

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	ms.Registry.MustRegister(srvMetrics)

	return MetricResults{
		MetricRegistry: &ms,
		MetricServer:   srvMetrics,
	}, nil
}
