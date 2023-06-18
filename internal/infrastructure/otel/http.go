package otel

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

func WrapAndReturn(httpTransport http.RoundTripper) *otelhttp.Transport {
	return otelhttp.NewTransport(httpTransport)
}
