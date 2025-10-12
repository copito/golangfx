package handler

import (
	"net/http"

	"github.com/copito/runner/src/internal/entities"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler is an http.Handler that copies its request body back to the response.
type MetricsHandler struct {
	metricStore *entities.MetricStore
}

// NewSwaggerFileHandler builds a new SwaggerFileHandler.
func NewMetricsHandler(metricStore *entities.MetricStore) *MetricsHandler {
	return &MetricsHandler{metricStore: metricStore}
}

func (h *MetricsHandler) Pattern() runtime.Pattern {
	// "/metrics"
	pattern, err := runtime.NewPattern(
		validPatternVersion,
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

func (h *MetricsHandler) Method() string {
	return "GET"
}

// ServeHTTP handles an HTTP request to the /openapi/* endpoint.
func (h *MetricsHandler) ServeHTTP() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Add prometheus metrics as route for http
		promhttp.HandlerFor(h.metricStore.Registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}
}
