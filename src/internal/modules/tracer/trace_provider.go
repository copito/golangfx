package tracer

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	otlpg "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/entities"
	"github.com/copito/runner/src/internal/modules/config"
)

type TraceProviderParams struct {
	fx.In
	Lifecycle      fx.Lifecycle
	Logger         *slog.Logger
	ConfigProvider config.ConfigProvider
	MetricStore    *entities.MetricStore
}

type TraceProviderResult struct {
	fx.Out
	TraceProvider *sdktrace.TracerProvider
}

func NewTraceProvider(params TraceProviderParams) (TraceProviderResult, error) {
	params.Logger.Info("setting up trace provider (with net)")
	config := params.ConfigProvider.Get()
	ctx := context.Background()

	var exporter sdktrace.SpanExporter
	switch config.OpenTelemetry.Type {
	case entities.OpenTelemetryTypeSTDOUT:
		// Set up OTLP tracing (stdout for debug)
		newExporter, err := stdout.New(stdout.WithPrettyPrint())
		if err != nil {
			params.Logger.Error("failed to init exporter", "err", slog.Any("err", err))
			panic(err)
		}
		exporter = newExporter
	case entities.OpenTelemetryTypeGRPC:
		// Set up OTLP tracing (endpoint)
		newExporter, err := otlpg.New(
			ctx,
			otlpg.WithEndpoint(config.OpenTelemetry.CollectorEndpoint),
			otlpg.WithInsecure(),
		)
		if err != nil {
			params.Logger.Error("failed to init exporter", "err", slog.Any("err", err))
			panic(err)
		}
		exporter = newExporter
	case entities.OpenTelemetryTypeHTTP:
		// Set up OTLP tracing (endpoint)
		newExporter, err := otlp.New(
			ctx,
			otlp.WithEndpoint(config.OpenTelemetry.CollectorEndpoint),
		)
		if err != nil {
			params.Logger.Error("failed to init exporter", "err", slog.Any("err", err))
			panic(err)
		}
		exporter = newExporter
	case entities.OpenTelemetryTypeDisabled:
		params.Logger.Error("OpenTelemetry is disabled")
		return TraceProviderResult{}, nil
	default:
		params.Logger.Error("invalid OpenTelemetry is disabled")
		return TraceProviderResult{}, fmt.Errorf("invalid OpenTelemetry type: %s", config.OpenTelemetry.Type)
	}

	// Setup OTLP tracing Sampler
	var sampler sdktrace.Sampler
	if config.Backend.Environment == entities.BackendEnvironmentLocal {
		sampler = sdktrace.AlwaysSample()
	} else {
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(config.OpenTelemetry.SamplingRate))
	}

	// DEFINE THE RESOURCE WITH THE SERVICE NAME
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.Global.Service),
			semconv.DeploymentEnvironmentKey.String(string(config.Backend.Environment)),
		),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
	)
	if err != nil {
		params.Logger.Error("failed to create otel resource", slog.Any("err", err))
		panic(err)
	}

	// Setup OTLP tracing
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// Register the tracer provider with Fx lifecycle
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			params.Logger.Info("starting trace provider...")
			return nil
		},
		OnStop: func(context.Context) error {
			params.Logger.Info("shutting down trace - flushing traces before shutdown...")
			return tp.Shutdown(ctx)
		},
	})

	return TraceProviderResult{
		TraceProvider: tp,
	}, nil
}
