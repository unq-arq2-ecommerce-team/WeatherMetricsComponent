package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/dto"
	"net/http"
)

// GetCityLastWeekTemperatureAverageHandler
// @Summary      Endpoint get city last week temperature average
// @Description  get city last week temperature average
// @Param city path string true "City" example("Quilmes")
// @Tags         Weather
// @Produce json
// @Success 200 {object} domain.AverageTemperature
// @Success 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router       /api/v1/weather/city/{city}/temperature/last/week [get]
func GetCityLastWeekTemperatureAverageHandler(logger domain.Logger, getCityLastWeekTemperatureAverageQuery *application.GetCityLastWeekTemperatureAverageQuery) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.WithRequestId(c)
		cityParam, _ := c.Params.Get("city")
		if cityParam == "" {
			log.WithFields(domain.LoggerFields{"cityParam": cityParam}).Errorf("city param is empty")
			c.JSON(http.StatusBadRequest, dto.NewErrorMessage("bad request", "empty city path param"))
			return
		}

		avgTemp, err := getCityLastWeekTemperatureAverageQuery.Do(c.Request.Context(), cityParam)
		if err != nil {
			switch err.(type) {
			case domain.AverageTemperatureNotFoundErr:
				c.JSON(http.StatusNotFound, dto.NewErrorMessage("last week average temperature of city not found", err.Error()))
			default:
				c.JSON(http.StatusInternalServerError, dto.NewErrorMessage("internal server error", err.Error()))
			}
			return
		}
		c.JSON(http.StatusOK, avgTemp)

	}
}
