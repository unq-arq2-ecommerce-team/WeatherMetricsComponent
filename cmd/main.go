package main

import (
	c "github.com/patrickmn/go-cache"
	"github.com/sony/gobreaker"
	app "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/application"
	infra "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache"
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

	var settings gobreaker.Settings
	settings.Name = "WeatherMetricsComponentBreaker"
	settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= uint32(conf.CircuitBreaker.MinRequests) && failureRatio >= conf.CircuitBreaker.FailuresRatio
	}
	settings.Timeout = conf.CircuitBreaker.Timeout
	settings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		if to == gobreaker.StateOpen {
			logger.Errorf("Circuit breaker named %s is open", name)
		}
		if from == gobreaker.StateOpen && to == gobreaker.StateHalfOpen {
			logger.Infof("Circuit breaker named %s from Open to  Half-open", name)
		}
		if from == gobreaker.StateHalfOpen && to == gobreaker.StateClosed {
			logger.Infof("Circuit breaker named %s from Half-open to Closed", name)
		}
	}
	cb := gobreaker.NewCircuitBreaker(settings)
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
