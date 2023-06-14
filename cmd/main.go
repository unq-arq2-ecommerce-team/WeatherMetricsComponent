package main

import (
	c "github.com/patrickmn/go-cache"
	app "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/application"
	infra "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache"
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
	})

	cacheTables := initCacheTables(conf.Weather.CurrentTemp.Cache, conf.Weather.AvgTemp.Cache)
	cacheClient := cache.NewLocalMemoryCacheClient(logger, cacheTables)

	// repositories
	baseWeatherRepo := http.NewWeatherRepository(logger, http.NewClient(conf.Weather.HttpConfig), conf.Weather)
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
