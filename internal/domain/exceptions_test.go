package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_WeatherNotFoundError(t *testing.T) {
	e1 := WeatherNotFoundError{
		City: "Quilmes",
	}
	e2 := WeatherNotFoundError{
		City: "Miamee",
	}
	assert.Equal(t, `weather with city Quilmes not found`, e1.Error())
	assert.Equal(t, `weather with city Miamee not found`, e2.Error())
}

func Test_AverageTemperatureNotFoundError(t *testing.T) {
	dateFrom1 := time.Date(2023, 10, 1, 12, 12, 12, 0, time.UTC)
	dateTo1 := time.Date(2023, 10, 5, 12, 12, 12, 0, time.UTC)
	dateFrom2 := time.Date(2021, 10, 1, 12, 12, 12, 0, time.UTC)
	dateTo2 := time.Date(2022, 10, 1, 12, 12, 12, 0, time.UTC)
	e1 := AverageTemperatureNotFoundErr{
		City:     "Quilmes",
		DateFrom: dateFrom1,
		DateTo:   dateTo1,
	}
	e2 := AverageTemperatureNotFoundErr{
		City:     "Miamee",
		DateFrom: dateFrom2,
		DateTo:   dateTo2,
	}
	assert.Equal(t, `average temperature with city Quilmes between date range [2023-10-01 12:12:12 +0000 UTC, 2023-10-05 12:12:12 +0000 UTC] not found`, e1.Error())
	assert.Equal(t, `average temperature with city Miamee between date range [2021-10-01 12:12:12 +0000 UTC, 2022-10-01 12:12:12 +0000 UTC] not found`, e2.Error())
}
