package domain

import "context"

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Timestamp   string  `json:"timestamp"`
}

type WeatherRepository interface {
	FindCurrentTemperatureByCity(ctx context.Context, city string) (*Weather, error)
}
