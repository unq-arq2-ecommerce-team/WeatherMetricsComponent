package application

import (
	"context"

	"github.com/sony/gobreaker"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
)

type GetCityLastDayTemperatureAverageQuery struct {
	weatherRepo domain.WeatherRepository
	cb          *gobreaker.CircuitBreaker
}

func NewGetCityLastDayTemperatureAverageQuery(weatherRepo domain.WeatherRepository, cb *gobreaker.CircuitBreaker) *GetCityLastDayTemperatureAverageQuery {
	return &GetCityLastDayTemperatureAverageQuery{
		weatherRepo: weatherRepo,
		cb:          cb,
	}
}

func (q *GetCityLastDayTemperatureAverageQuery) Do(ctx context.Context, city string) (*domain.AverageTemperature, error) {
	dateFrom, dateTo := domain.GetLastDayDates()
	avgTemp, err := q.cb.Execute(func() (interface{}, error) {
		return q.weatherRepo.GetAverageTemperatureByCityAndDateRange(ctx, city, dateFrom, dateTo)
	})
	return domain.ParseOrNil[domain.AverageTemperature](avgTemp), err
}
