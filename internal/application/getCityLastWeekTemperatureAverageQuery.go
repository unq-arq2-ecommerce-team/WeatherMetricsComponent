package application

import (
	"context"

	"github.com/sony/gobreaker"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
)

type GetCityLastWeekTemperatureAverageQuery struct {
	weatherRepo domain.WeatherRepository
	cb          *gobreaker.CircuitBreaker
}

func NewGetCityLastWeekTemperatureAverageQuery(weatherRepo domain.WeatherRepository, cb *gobreaker.CircuitBreaker) *GetCityLastWeekTemperatureAverageQuery {
	return &GetCityLastWeekTemperatureAverageQuery{
		weatherRepo: weatherRepo,
		cb:          cb,
	}
}

func (q *GetCityLastWeekTemperatureAverageQuery) Do(ctx context.Context, city string) (*domain.AverageTemperature, error) {
	dateFrom, dateTo := domain.GetLastWeekDates()

	avgTemp, err := q.cb.Execute(func() (interface{}, error) {
		return q.weatherRepo.GetAverageTemperatureByCityAndDateRange(ctx, city, dateFrom, dateTo)
	})
	return avgTemp.(*domain.AverageTemperature), err
}
