package application

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
)

type FindCityCurrentTemperatureQuery struct {
	weatherRepo domain.WeatherRepository
}

func NewFindCityCurrentTemperatureQuery(weatherRepo domain.WeatherRepository) *FindCityCurrentTemperatureQuery {
	return &FindCityCurrentTemperatureQuery{
		weatherRepo: weatherRepo,
	}
}

func (q *FindCityCurrentTemperatureQuery) Do(ctx context.Context, city string) (*domain.Weather, error) {
	return q.weatherRepo.FindCurrentTemperatureByCity(ctx, city)
}
