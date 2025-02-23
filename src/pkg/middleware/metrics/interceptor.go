package metrics

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	"github.com/copito/runner/src/internal/entities"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"google.golang.org/grpc"
)

type MetricInterceptor struct {
	logger *slog.Logger

	promLabels func(ctx context.Context) prometheus.Labels
	SrvMetrics *grpcprom.ServerMetrics
}

func NewMetricInterceptor(logger *slog.Logger, config *entities.Config, ms *entities.MetricStore) *MetricInterceptor {
	// Setup Metrics
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	ms.Registry.MustRegister(srvMetrics)
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		span := trace.SpanContextFromContext(ctx)
		if span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	var exporter sdktrace.SpanExporter
	if config.Backend.OpenTelemetry.Type == "STDOUT" {
		// Set up OTLP tracing (stdout for debug)
		newExporter, err := stdout.New(stdout.WithPrettyPrint())
		if err != nil {
			logger.Error("failed to init exporter", "err", slog.Any("err", err))
			panic(err)
		}

		exporter = newExporter
	} else if config.Backend.OpenTelemetry.Type == "HTTP" {
		// Set up OTLP tracing (endpoint)
		newExporter, err := otlp.New(context.Background(), otlp.WithEndpoint(config.Backend.OpenTelemetry.CollectorEndpoint))
		if err != nil {
			logger.Error("failed to init exporter", "err", slog.Any("err", err))
			panic(err)
		}
		exporter = newExporter
		// } else if config.Backend.OpenTelemetry.Type == "GRPC" {
		// 	// Set up OTLP tracing (endpoint)
		// 	newExporter, err := otlp.New(context.Background(), otlp.WithEndpoint(config.Backend.OpenTelemetry.CollectorEndpoint))
		// 	if err != nil {
		// 		logger.Error("failed to init exporter", "err", slog.Any("err", err))
		// 		panic(err)
		// 	}
		// 	exporter = newExporter
	} else {
		logger.Error("unsupported open telemetry type", slog.String("type", config.Backend.OpenTelemetry.Type))
		panic("invalid OpenTelemetry type")
	}

	// Setup OTLP tracing Sampler
	var sampler sdktrace.Sampler
	if config.Backend.Environment == "local" {
		sampler = sdktrace.AlwaysSample()
	} else {
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.01))
	}

	// Setup OTLP tracing
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	// _ = exporter.Shutdown(context.Background())
	// defer func() { _ = exporter.Shutdown(context.Background()) }()

	return &MetricInterceptor{
		logger:     logger,
		promLabels: exemplarFromContext,
		SrvMetrics: srvMetrics,
	}
}

func (l MetricInterceptor) BuildUnaryInterceptor() grpc.UnaryServerInterceptor {
	return l.SrvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(l.promLabels))
}

func (l MetricInterceptor) BuildStreamInterceptor() grpc.StreamServerInterceptor {
	return l.SrvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(l.promLabels))
}
