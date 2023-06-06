package http

import (
	"context"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
	"net/http"
)

func NewClient() *http.Client {
	return cleanhttp.DefaultPooledClient()
}

func IsStatusCode2XX(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func NewRequestWithContextWithNoBody(ctx context.Context, httpMethod, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, nil)
	req.Header.Add(logger.RequestIdHeaderKey(), logger.GetRequestId(ctx))
	return req, err
}
