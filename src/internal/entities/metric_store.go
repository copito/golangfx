package entities

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricStore struct {
	Registry *prometheus.Registry

	// Total number of panics
	PanicsTotal prometheus.Counter
}
