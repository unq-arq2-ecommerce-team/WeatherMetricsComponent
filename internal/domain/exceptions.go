package domain

import (
	"fmt"
	"time"
)

type WeatherNotFoundError struct {
	City string
}

func (e WeatherNotFoundError) Error() string {
	return fmt.Sprintf("weather with city %s not found", e.City)
}

type AverageTemperatureNotFoundErr struct {
	City     string
	DateFrom time.Time
	DateTo   time.Time
}

func NewAverageTemperatureNotFoundError(city string, dateFrom, dateTo time.Time) AverageTemperatureNotFoundErr {
	return AverageTemperatureNotFoundErr{
		City:     city,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
}

func (e AverageTemperatureNotFoundErr) Error() string {
	return fmt.Sprintf("average temperature with city %s between date range [%s, %s] not found", e.City, e.DateFrom, e.DateTo)
}
