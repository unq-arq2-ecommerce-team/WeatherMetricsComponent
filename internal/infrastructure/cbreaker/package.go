package cbreaker

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sony/gobreaker"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
)

func NewCircuitBreaker(logger domain.Logger, conf config.CircuitBreakerConfig) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(getSettings(logger, conf))
}

func getSettings(logger domain.Logger, conf config.CircuitBreakerConfig) gobreaker.Settings {
	log := logger.WithFields(domain.LoggerFields{"loggerFrom": "circuit breaker"})
	var settings gobreaker.Settings
	settings.Name = "WeatherMetricsComponentBreaker"
	settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= uint32(conf.MinRequests) && failureRatio >= conf.FailuresRatio
	}
	settings.Timeout = conf.Timeout
	settings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		if to == gobreaker.StateOpen {
			log.Errorf("Circuit breaker named %s is open", name)
		}
		if from == gobreaker.StateOpen && to == gobreaker.StateHalfOpen {
			log.Infof("Circuit breaker named %s from Open to Half-open", name)
		}
		if from == gobreaker.StateHalfOpen && to == gobreaker.StateClosed {
			log.Infof("Circuit breaker named %s from Half-open to Closed", name)
		}
		circuitBreakerStatusChange(to)
	}
	return settings
}

var CircuitBreakerStatus = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "circuit_breaker_status",
	Help: "Current status of the circuit breaker (closed, half_open, open)",
})

func circuitBreakerStatusChange(newState gobreaker.State) {
	CircuitBreakerStatus.Set(float64(newState))
}
