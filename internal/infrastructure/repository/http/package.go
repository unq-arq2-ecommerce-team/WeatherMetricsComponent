package http

import (
	"context"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

const deltaRetryWait = 2 * time.Second

func NewDefaultClient() *http.Client {
	client := cleanhttp.DefaultPooledClient()
	client.Transport = otelhttp.NewTransport(client.Transport)
	return cleanhttp.DefaultPooledClient()
}

func NewClient(logger domain.Logger, httpConfig config.HttpConfig) *http.Client {
	httpClient := NewDefaultClient()
	httpClient.Timeout = httpConfig.Timeout

	retryableClient := retryablehttp.NewClient()
	retryableClient.Logger = logger.WithFields(domain.LoggerFields{"loggerFrom": "http.retryableClient"})
	retryableClient.HTTPClient = httpClient
	retryableClient.RetryMax = httpConfig.Retries
	retryableClient.RetryWaitMin = httpConfig.RetryWait
	retryableClient.RetryWaitMax = httpConfig.RetryWait + deltaRetryWait
	return retryableClient.StandardClient()
}

func IsStatusCode2XX(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func NewRequestWithContextWithNoBody(ctx context.Context, httpMethod, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, nil)
	req.Header.Add(logger.RequestIdHeaderKey(), logger.GetRequestId(ctx))
	return req, err
}
