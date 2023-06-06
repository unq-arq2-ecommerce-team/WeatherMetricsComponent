package main

import (
	app "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/application"
	infra "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure"
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

	// repositories
	weatherRepo := http.NewWeatherRepository(logger, http.NewClient(), conf.Weather)

	// use cases
	findCityCurrentTemperatureQuery := app.NewFindCityCurrentTemperatureQuery(weatherRepo)

	application := infra.NewGinApplication(conf, logger, findCityCurrentTemperatureQuery)
	logger.Fatal(application.Run())
}
