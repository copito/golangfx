package modules

import (
	"context"
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
)

type TraceProviderParams struct {
	fx.In
	Lifecycle   fx.Lifecycle
	Logger      *slog.Logger
	Config      *entities.Config
	MetricStore *entities.MetricStore
}

type TraceProviderResult struct {
	fx.Out
	TraceProvider *sdktrace.TracerProvider
}

func NewTraceProvider(params TraceProviderParams) (TraceProviderResult, error) {
	params.Logger.Info("setting up trace provider (with net)")
	ctx := context.Background()

	var exporter sdktrace.SpanExporter
	switch params.Config.OpenTelemetry.Type {
	case "STDOUT":
		// Set up OTLP tracing (stdout for debug)
		newExporter, err := stdout.New(stdout.WithPrettyPrint())
		if err != nil {
			params.Logger.Error("failed to init exporter", "err", slog.Any("err", err))
			panic(err)
		}
		exporter = newExporter
	case "GRPC":
		// Set up OTLP tracing (endpoint)
		newExporter, err := otlpg.New(
			ctx,
			otlpg.WithEndpoint(params.Config.OpenTelemetry.CollectorEndpoint),
			otlpg.WithInsecure(),
		)
		if err != nil {
			params.Logger.Error("failed to init exporter", "err", slog.Any("err", err))
			panic(err)
		}
		exporter = newExporter
	case "HTTP":
		// Set up OTLP tracing (endpoint)
		newExporter, err := otlp.New(
			ctx,
			otlp.WithEndpoint(params.Config.OpenTelemetry.CollectorEndpoint),
		)
		if err != nil {
			params.Logger.Error("failed to init exporter", "err", slog.Any("err", err))
			panic(err)
		}
		exporter = newExporter
	case "DISABLED":
		params.Logger.Error("OpenTelemetry is disabled")
	default:
		params.Logger.Error("invalid OpenTelemetry is disabled")
	}

	// Setup OTLP tracing Sampler
	var sampler sdktrace.Sampler
	if params.Config.Backend.Environment == "local" {
		sampler = sdktrace.AlwaysSample()
	} else {
		sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(params.Config.OpenTelemetry.SamplingRate))
	}

	// DEFINE THE RESOURCE WITH THE SERVICE NAME
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(params.Config.Global.Service),
			semconv.DeploymentEnvironmentKey.String(params.Config.Backend.Environment),
		),
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

var TraceProviderModule = fx.Provide(NewTraceProvider)
