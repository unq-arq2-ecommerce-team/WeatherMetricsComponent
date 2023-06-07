package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
)

type GetCityLastWeekTemperatureAverageQuery struct {
	weatherRepo domain.WeatherRepository
}

func NewGetCityLastWeekTemperatureAverageQuery(weatherRepo domain.WeatherRepository) *GetCityLastWeekTemperatureAverageQuery {
	return &GetCityLastWeekTemperatureAverageQuery{
		weatherRepo: weatherRepo,
	}
}

func (q *GetCityLastWeekTemperatureAverageQuery) Do(ctx context.Context, city string) (*domain.AverageTemperature, error) {
	dateFrom, dateTo := domain.GetLastWeekDates()
	return q.weatherRepo.GetAverageTemperatureByCityAndDateRange(ctx, city, dateFrom, dateTo)
}
