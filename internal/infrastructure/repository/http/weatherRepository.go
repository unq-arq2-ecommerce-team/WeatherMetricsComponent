package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
	"io"
	"net/http"
	"strings"
	"time"
)

type weatherRepository struct {
	logger                domain.Logger
	client                *http.Client
	getCurrentTempUrl     string
	getAvgTempLastDayUrl  string
	getAvgTempLastWeekUrl string
}

func NewWeatherRepository(logger domain.Logger, client *http.Client, weatherConfig config.WeatherEndpoint) domain.WeatherRepository {
	return &weatherRepository{
		logger:                logger.WithFields(domain.LoggerFields{"repository.http": "weatherRepository"}),
		client:                client,
		getCurrentTempUrl:     weatherConfig.CurrentTempUrl,
		getAvgTempLastDayUrl:  weatherConfig.AvgTempLastDayUrl,
		getAvgTempLastWeekUrl: weatherConfig.AvgTempLastWeekUrl,
	}
}

func (r weatherRepository) FindCurrentTemperatureByCity(ctx context.Context, city string) (*domain.Weather, error) {
	url := strings.Replace(r.getCurrentTempUrl, "{city}", city, -1)

	log := r.logger.WithRequestId(ctx).WithFields(loggerPkg.Fields{"method": "findCurrentTemperatureByCity", "url": url})

	req, err := NewRequestWithContextWithNoBody(ctx, http.MethodGet, url)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when create request")
		return nil, err
	}
	sw := time.Now()
	res, err := r.client.Do(req)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Error("http error do request")
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	log.Debugf("request finished in %s", time.Since(sw))

	log = log.WithFields(domain.LoggerFields{"statusCode": res.StatusCode})

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Error("an error has occurred sending http request")
		return nil, err
	}
	log = log.WithFields(domain.LoggerFields{"responseBody": string(rawBody)})

	switch statusCode := res.StatusCode; {
	case IsStatusCode2XX(statusCode):
		var weather domain.Weather
		if err := json.Unmarshal(rawBody, &weather); err != nil {
			log.WithFields(domain.LoggerFields{"error": err}).Errorf("error decoding res body")
			return nil, fmt.Errorf("weather repository error with status code %v and url %s", statusCode, url)
		}
		log.Infof("successful find current temperature with city %s", city)
		return &weather, nil
	case statusCode == http.StatusNotFound:
		return nil, domain.WeatherNotFoundError{City: city}
	default:
		return nil, fmt.Errorf("weather repository error with status code %v and url %s", statusCode, url)
	}
}
