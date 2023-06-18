package otel

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitOtelTrace(ctx context.Context, logger domain.Logger, otelConf config.OtelConfig, isEnabled bool) func() {
	if isEnabled {
		cleanupFn := initTracerAuto(logger, otelConf, config.OtlServiceName, config.ServiceName)
		return func() {
			err := cleanupFn(ctx)
			if err != nil {
				logger.WithFields(domain.LoggerFields{"error": err}).Errorf("some error found when clean up applied")
			}
		}
	}
	return func() {}
}

func initTracerAuto(baseLogger domain.Logger, conf config.OtelConfig, serviceName, appName string) func(context.Context) error {
	logger := baseLogger.WithFields(domain.LoggerFields{"loggingFrom": "initTracerAuto"})
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(conf.URL),
		),
	)

	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Fatalf("Could not set exporter")
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("application", appName),
		),
	)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Fatalf("Could not set resources: %s", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resources),
		),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return exporter.Shutdown
}
