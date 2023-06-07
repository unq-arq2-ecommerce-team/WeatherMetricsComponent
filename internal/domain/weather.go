package domain

import (
	"context"
	"time"
)

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Timestamp   string  `json:"timestamp"`
}

type AverageTemperature struct {
	City           string    `json:"city" bson:"_id"`
	AvgTemperature float64   `json:"avgTemperature" bson:"avgTemperature"`
	DateFrom       time.Time `json:"dateFrom"`
	DateTo         time.Time `json:"dateTo"`
	DaysBetween    float64   `json:"daysBetween"`
}

func (w Weather) String() string {
	return ParseStruct(w)
}

func (a *AverageTemperature) String() string {
	return ParseStruct(a)
}

type WeatherRepository interface {
	FindCurrentTemperatureByCity(ctx context.Context, city string) (*Weather, error)
	GetAverageTemperatureByCityAndDateRange(ctx context.Context, city string, dateFrom, dateTo time.Time) (*AverageTemperature, error)
}
