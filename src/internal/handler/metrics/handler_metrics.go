package metrics

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/internal/handler/common"
)

type MetricsHandler interface {
	common.HttpHandlerInterface
}

// MetricsHandler is an http.Handler that copies its request body back to the response.
type metricsHandler struct {
	metricStore *entities.MetricStore
}

// NewSwaggerFileHandler builds a new SwaggerFileHandler.
func NewMetricsHandler(metricStore *entities.MetricStore) *metricsHandler {
	return &metricsHandler{metricStore: metricStore}
}

func (h *metricsHandler) Pattern() runtime.Pattern {
	// "/metrics"
	pattern, err := runtime.NewPattern(
		common.ValidPatternVersion,
		[]int{
			int(utilities.OpLitPush), 0, // runtime.OpLitPush → Push the literal "metrics" (matches /metrics exactly).
		},
		[]string{"metrics"},
		"", // no verb (gRPC routing suffix)
	)
	if err != nil {
		panic("error registering pattern for swagger file handler")
	}
	return pattern
}

func (h *metricsHandler) Method() string {
	return "GET"
}

// ServeHTTP handles an HTTP request to the /openapi/* endpoint.
func (h *metricsHandler) ServeHTTP() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Add prometheus metrics as route for http
		promhttp.HandlerFor(h.metricStore.Registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}
}
