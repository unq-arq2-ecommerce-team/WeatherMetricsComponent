package main

import (
	app "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/application"
	infra "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache"
	localCache "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache/local"
	// redisCache "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache/redis"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cbreaker"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/repository/http"
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

	// cache client
	localCacheClient := localCache.NewLocalMemoryCacheClient(logger, conf.LocalCache)
	// redisCacheClient := redisCache.NewCacheClient(logger, conf.Redis)

	// repositories
	baseWeatherRepo := http.NewWeatherRepository(logger, http.NewClient(logger, conf.Weather.HttpConfig), conf.Weather)
	cacheWeatherRepo := cache.NewWeatherRepository(logger, localCacheClient, baseWeatherRepo)

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
