package main

import (
	"context"
	c "github.com/patrickmn/go-cache"
	app "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/application"
	infra "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cbreaker"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/repository/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"log"
)

func main() {
	conf := config.LoadConfig()
	logger := loggerPkg.New(&loggerPkg.Config{
		ServiceName:     config.ServiceName,
		EnvironmentName: conf.Environment,
		LogLevel:        conf.LogLevel,
		LogFormat:       loggerPkg.JsonFormat,
		LokiHost:        conf.LokiHost,
	})

	cacheTables := initCacheTables(conf.Weather.CurrentTemp.Cache, conf.Weather.AvgTemp.Cache)
	cacheClient := cache.NewLocalMemoryCacheClient(logger, cacheTables)

	// OTEL
	cleanup := initTracerAuto()
	defer cleanup(context.Background())

	// repositories
	baseWeatherRepo := http.NewWeatherRepository(logger, http.NewClient(logger, conf.Weather.HttpConfig), conf.Weather)
	cacheWeatherRepo := cache.NewWeatherRepository(logger, cacheClient, baseWeatherRepo)

	//circuit breaker
	cb := cbreaker.NewCircuitBreaker(logger, conf.CircuitBreaker)

	// use cases
	findCityCurrentTemperatureQuery := app.NewFindCityCurrentTemperatureQuery(cacheWeatherRepo, cb)
	getCityLastDayTemperatureAverageQuery := app.NewGetCityLastDayTemperatureAverageQuery(cacheWeatherRepo, cb)
	getCityLastWeekTemperatureAverageQuery := app.NewGetCityLastWeekTemperatureAverageQuery(cacheWeatherRepo, cb)

	application := infra.NewGinApplication(
		conf,
		logger,
		findCityCurrentTemperatureQuery,
		getCityLastDayTemperatureAverageQuery,
		getCityLastWeekTemperatureAverageQuery,
	)
	logger.Fatal(application.Run())
}

func initCacheTables(currentTemp, avgTemp config.CacheConfig) map[string]*c.Cache {
	return map[string]*c.Cache{
		cache.TableCurrentTemperature: c.New(currentTemp.DefaultExp, currentTemp.PurgesExp),
		cache.TableAvgTemperature:     c.New(avgTemp.DefaultExp, avgTemp.PurgesExp),
	}
}

func initTracerAuto() func(ctx context.Context) error {

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("otel-collector:4317"),
		),
	)

	if err != nil {
		log.Fatal("Could not set exporter: ", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "weather-metrics"),
			attribute.String("application", "WeatherMetricsComponent"),
		),
	)
	if err != nil {
		log.Fatal("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resources),
		),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return exporter.Shutdown
}
