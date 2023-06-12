package cache

import (
	"context"
	"fmt"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"time"
)

const (
	layoutKeyTime           = time.RFC3339
	TableCurrentTemperature = "temperature:current"
	TableAvgTemperature     = "temperature:average"
)

type cacheWeatherRepository struct {
	logger            domain.Logger
	cacheClient       MemoryCacheClient
	weatherRepository domain.WeatherRepository
}

func NewWeatherRepository(logger domain.Logger, cacheClient MemoryCacheClient, weatherRepository domain.WeatherRepository) domain.WeatherRepository {
	return &cacheWeatherRepository{
		logger:            logger.WithFields(domain.LoggerFields{"repository.cache": "cacheWeatherRepository"}),
		cacheClient:       cacheClient,
		weatherRepository: weatherRepository,
	}
}

func (c cacheWeatherRepository) FindCurrentTemperatureByCity(ctx context.Context, city string) (*domain.Weather, error) {
	log := c.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"method": "FindCurrentTemperatureByCity", "table": TableCurrentTemperature, "city": city})

	// get cache and returns if was found
	cachedWeatherRaw, found, err := c.cacheClient.Get(ctx, TableCurrentTemperature, city)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("cache table not found")
	}
	if found {
		log.Debug("successful get FindCurrentTemperatureByCity from cache")
		return ParseWeather(log, cachedWeatherRaw)
	}

	// original repository method
	weather, err := c.weatherRepository.FindCurrentTemperatureByCity(ctx, city)
	if err != nil {
		return weather, err
	}

	// save in cache
	expirationTime := getExpiresTimeInCurrentTemperature(weather.Timestamp)
	log = log.WithFields(domain.LoggerFields{"expirationTime": expirationTime.String()})
	if err := c.cacheClient.Save(ctx, TableCurrentTemperature, city, weather, expirationTime); err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error save cache client")
	} else {
		log.Debug("successful FindCurrentTemperatureByCity cached")
	}
	return weather, nil
}

func (c cacheWeatherRepository) GetAverageTemperatureByCityAndDateRange(ctx context.Context, city string, dateFrom, dateTo time.Time) (*domain.AverageTemperature, error) {
	cacheKey := getAvgTempCacheKey(city, dateFrom, dateTo)
	log := c.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"method": "GetAverageTemperatureByCityAndDateRange", "table": TableAvgTemperature, "cacheKey": cacheKey, "city": city, "dateFrom": dateFrom, "dateTo": dateTo})

	// get cache and returns if was found
	cachedAvgTempRaw, found, err := c.cacheClient.Get(ctx, TableAvgTemperature, cacheKey)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("cache table not found")
	}
	if found {
		log.Debug("successful get GetAverageTemperatureByCityAndDateRange from cache")
		return ParseAvgTemp(log, cachedAvgTempRaw)
	}
	// original repository method
	avgTemp, err := c.weatherRepository.GetAverageTemperatureByCityAndDateRange(ctx, city, dateFrom, dateTo)
	if err != nil {
		return avgTemp, err
	}
	// save in cache
	expirationTime := getExpiresTimeInAverageTemperature()
	log = log.WithFields(domain.LoggerFields{"expirationTime": expirationTime.String()})
	if err := c.cacheClient.Save(ctx, TableAvgTemperature, cacheKey, avgTemp, expirationTime); err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error save cache client with table %s key %s", TableAvgTemperature, city)
	} else {
		log.Debug("successful GetAverageTemperatureByCityAndDateRange cached")
	}
	return avgTemp, nil
}

// getExpiresTimeInCurrentTemperature : is always difference between now date and next hour
func getExpiresTimeInCurrentTemperature(timestamp time.Time) time.Duration {
	return timestamp.Add(time.Hour).Sub(time.Now().UTC())
}

// getExpiresTimeInAverageTemperature : is always difference between tomorrow date and date now (because it's always refreshed every day)
func getExpiresTimeInAverageTemperature() time.Duration {
	timeNow := time.Now().UTC()
	return domain.GetFollowingDay(timeNow).Sub(timeNow)
}

// getAvgTempCacheKey returns a string compose by city and dateFrom day and dateTo day
func getAvgTempCacheKey(city string, dateFrom, dateTo time.Time) string {
	return fmt.Sprintf("%s_%s_%s", city, dateFrom.Format(layoutKeyTime), dateTo.Format(layoutKeyTime))
}

func ParseWeather(log domain.Logger, data interface{}) (*domain.Weather, error) {
	weather, ok := data.(*domain.Weather)
	if !ok {
		err := fmt.Errorf("raw data cannot be parsed to *domain.Weather")
		log.WithFields(domain.LoggerFields{"error": err}).Error("error parsing data")
		return nil, err
	}
	log.Info("successful get weather from cache")
	return weather, nil
}

func ParseAvgTemp(log domain.Logger, data interface{}) (*domain.AverageTemperature, error) {
	avgTemp, ok := data.(*domain.AverageTemperature)
	if !ok {
		err := fmt.Errorf("raw data cannot be parsed to *domain.AverageTemperature")
		log.WithFields(domain.LoggerFields{"error": err}).Error("error parsing data")
		return nil, err
	}
	log.Info("successful get average temperature from cache")
	return avgTemp, nil
}
