package application

import (
	"context"

	"github.com/sony/gobreaker"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
)

type FindCityCurrentTemperatureQuery struct {
	weatherRepo domain.WeatherRepository
	cb          *gobreaker.CircuitBreaker
}

func NewFindCityCurrentTemperatureQuery(weatherRepo domain.WeatherRepository, cb *gobreaker.CircuitBreaker) *FindCityCurrentTemperatureQuery {
	return &FindCityCurrentTemperatureQuery{
		weatherRepo: weatherRepo,
		cb:          cb,
	}
}

func (q *FindCityCurrentTemperatureQuery) Do(ctx context.Context, city string) (*domain.Weather, error) {

	weather, err := q.cb.Execute(func() (interface{}, error) {
		return q.weatherRepo.FindCurrentTemperatureByCity(ctx, city)
	})
	return weather.(*domain.Weather), err
}
