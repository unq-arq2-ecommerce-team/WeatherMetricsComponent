package dto

import (
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"time"
)

type TemperatureAverageSearch struct {
	DateFrom time.Time `json:"dateFrom" url:"dateFrom" layout:"2006-01-02T15:04:05.000Z" time_utc:"1"`
	DateTo   time.Time `json:"dateTo" url:"dateTo" layout:"2006-01-02T15:04:05.000Z" time_utc:"1"`
}

func (q TemperatureAverageSearch) String() string {
	return domain.ParseStruct(q)
}
