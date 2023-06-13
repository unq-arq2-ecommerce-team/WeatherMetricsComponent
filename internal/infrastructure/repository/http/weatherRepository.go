package http

import (
	"context"
	"encoding/json"
	"fmt"
	queryParam "github.com/google/go-querystring/query"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/dto"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
	"io"
	"net/http"
	"strings"
	"time"
)

type weatherRepository struct {
	logger            domain.Logger
	httpClient        *http.Client
	getCurrentTempUrl string
	getAvgTempUrl     string
}

func NewWeatherRepository(logger domain.Logger, httpClient *http.Client, weatherConfig config.WeatherEndpoint) domain.WeatherRepository {
	return &weatherRepository{
		logger:            logger.WithFields(domain.LoggerFields{"repository.http": "weatherRepository"}),
		httpClient:        httpClient,
		getCurrentTempUrl: weatherConfig.CurrentTemp.Url,
		getAvgTempUrl:     weatherConfig.AvgTemp.Url,
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
	res, err := r.httpClient.Do(req)
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
		log.Infof("successful find current temperature")
		return &weather, nil
	case statusCode == http.StatusNotFound:
		return nil, domain.WeatherNotFoundError{City: city}
	default:
		return nil, fmt.Errorf("weather repository error with status code %v and url %s", statusCode, url)
	}
}

func (r weatherRepository) GetAverageTemperatureByCityAndDateRange(ctx context.Context, city string, dateFrom, dateTo time.Time) (*domain.AverageTemperature, error) {
	url := strings.Replace(r.getAvgTempUrl, "{city}", city, -1)
	log := r.logger.WithRequestId(ctx).WithFields(loggerPkg.Fields{"method": "GetAverageTemperatureByCityAndDateRange", "url": url, "city": city, "dateFrom": dateFrom, "dateTo": dateTo})

	req, err := NewRequestWithContextWithNoBody(ctx, http.MethodGet, url)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when create request")
		return nil, err
	}
	values, err := queryParam.Values(dto.TemperatureAverageSearch{DateFrom: dateFrom, DateTo: dateTo})
	if err == nil {
		log.Debug("successful query param values %s", values.Encode())
		req.URL.RawQuery = values.Encode()
	}
	log = log.WithFields(domain.LoggerFields{"url": req.URL.String()})
	sw := time.Now()
	res, err := r.httpClient.Do(req)
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
		var avgTemp domain.AverageTemperature
		if err := json.Unmarshal(rawBody, &avgTemp); err != nil {
			log.WithFields(domain.LoggerFields{"error": err}).Errorf("error decoding res body")
			return nil, fmt.Errorf("weather repository error with status code %v and url %s", statusCode, url)
		}
		log.Infof("successful get temperature average")
		return &avgTemp, nil
	case statusCode == http.StatusNotFound:
		return nil, domain.NewAverageTemperatureNotFoundError(city, dateFrom, dateTo)
	default:
		return nil, fmt.Errorf("weather repository error with status code %v and url %s", statusCode, url)
	}
}
