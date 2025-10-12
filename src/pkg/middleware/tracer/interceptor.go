package tracer

import (
	"context"
	"log/slog"

	"github.com/copito/runner/src/internal/entities"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/trace"
)

type TracerInterceptor struct {
	logger *slog.Logger

	promLabels func(ctx context.Context) prometheus.Labels
	SrvMetrics *grpcprom.ServerMetrics
}

func NewTracerInterceptor(logger *slog.Logger, config *entities.Config, metricServer *grpcprom.ServerMetrics) TracerInterceptor {
	examplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	return TracerInterceptor{
		logger:     logger,
		promLabels: examplarFromContext,
		SrvMetrics: metricServer,
	}
}

func (l TracerInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	return l.SrvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(l.promLabels))
}

func (l TracerInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return l.SrvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(l.promLabels))
}
