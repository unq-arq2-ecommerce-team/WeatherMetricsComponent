package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
)

type GetCityLastDayTemperatureAverageQuery struct {
	weatherRepo domain.WeatherRepository
}

func NewGetCityLastDayTemperatureAverageQuery(weatherRepo domain.WeatherRepository) *GetCityLastDayTemperatureAverageQuery {
	return &GetCityLastDayTemperatureAverageQuery{
		weatherRepo: weatherRepo,
	}
}

func (q *GetCityLastDayTemperatureAverageQuery) Do(ctx context.Context, city string) (*domain.AverageTemperature, error) {
	dateFrom, dateTo := domain.GetLastDayDates()
	return q.weatherRepo.GetAverageTemperatureByCityAndDateRange(ctx, city, dateFrom, dateTo)
}
