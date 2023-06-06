package infrastructure

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerDocs "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/docs"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/handlers"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/middleware"
	"io"
	"net/http"
)

// Application
// @title Weather Metrics Component API
// @version 1.0
// @description api for final tp arq2
// @contact.name API SUPPORT
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
// @query.collection.format multi
type Application interface {
	Run() error
}

type ginApplication struct {
	logger                          domain.Logger
	config                          config.Config
	findCityCurrentTemperatureQuery *application.FindCityCurrentTemperatureQuery
}

func NewGinApplication(config config.Config, logger domain.Logger, findCityCurrentTemperatureQuery *application.FindCityCurrentTemperatureQuery) Application {
	return &ginApplication{
		logger:                          logger,
		config:                          config,
		findCityCurrentTemperatureQuery: findCityCurrentTemperatureQuery,
	}
}

func (app *ginApplication) Run() error {
	swaggerDocs.SwaggerInfo.Host = fmt.Sprintf("localhost:%v", app.config.Port)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router := gin.Default()
	router.GET("/", HealthCheck)

	routerApiV1 := router.Group("/api/v1")
	routerApiV1.Use(middleware.TracingRequestId())

	routerApiV1.GET("/weather/city/:city/temperature", handlers.FindCityCurrentTemperatureHandler(app.logger, app.findCityCurrentTemperatureQuery))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.logger.Infof("running http server on port %d", app.config.Port)
	return router.Run(fmt.Sprintf(":%v", app.config.Port))
}

// HealthCheck
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Health check
// @Accept */*
// @Produce json
// @Success 200 {object} HealthCheckRes
// @Router / [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckRes{Data: "Server is up and running"})
}

type HealthCheckRes struct {
	Data string `json:"data"`
}
