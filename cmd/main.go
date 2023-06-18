package main

import (
	"context"
	app "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/application"
	infra "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache"
	redisCache "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache/redis"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cbreaker"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/otel"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/repository/http"
)

func main() {
	conf := config.LoadConfig()
	isIntegrationEnv := conf.IsIntegrationEnv()

	logger := loggerPkg.New(&loggerPkg.Config{
		ServiceName:      config.ServiceName,
		EnvironmentName:  conf.Environment,
		IsIntegrationEnv: isIntegrationEnv,
		LogLevel:         conf.LogLevel,
		LogFormat:        loggerPkg.JsonFormat,
		LokiHost:         conf.LokiHost,
	})

	// OTEL
	cleanupFn := otel.InitOtelTrace(context.Background(), logger, conf.Otel, isIntegrationEnv)
	defer cleanupFn()

	// cache client
	// localCacheClient := localCache.NewLocalMemoryCacheClient(logger, conf.LocalCache)
	redisCacheClient := redisCache.NewCacheClient(logger, conf.Redis)

	// repositories
	conf.Weather.HttpConfig.OtelEnabled = isIntegrationEnv
	baseWeatherRepo := http.NewWeatherRepository(logger, http.NewClient(logger, conf.Weather.HttpConfig), conf.Weather)
	cacheWeatherRepo := cache.NewWeatherRepository(logger, redisCacheClient, baseWeatherRepo)

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
