package example

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/AndersonQ/talks/2021.11.13/example/config"
)

func SetUp(cfg config.Config) {
	otlpClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.OTelExporterEndpoint))

	otlpExporter, err := otlptrace.New(context.TODO(), otlpClient)
	if err != nil {
		cfg.Logger.Warn().Err(err).Msg("failed to create OTel exporter, disabling OTel")
		return
	}

	// TODO: comment
	res, err := resource.New(context.TODO(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String("develop"),
			semconv.DeploymentEnvironmentKey.String("local"),
			attribute.String("resource_attribute", "added as resource"),
		),
	)
	if err != nil {
		cfg.Logger.Warn().Err(err).Msg("failed to create otel sdk/resource")
	}

	// TODO: comment
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithBatcher(otlpExporter),
	)

	if cfg.Debug {
		cfg.Logger.Debug().Msg("adding stdout span processor")

		stdoutExporter, err := stdouttrace.New()
		if err != nil {
			cfg.Logger.Fatal().Err(err).Msg("failed to initialize stdouttrace export pipeline")
		}

		tracerProvider.RegisterSpanProcessor(
			trace.NewSimpleSpanProcessor(stdoutExporter))
	}

	// TODO comment
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)
}
