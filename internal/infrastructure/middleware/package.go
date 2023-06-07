package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
)

const headerRequestId = "system-request-id"

func TracingRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := logger.SetRequestId(c.Request.Context(), c.Request.Header.Get(headerRequestId))
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(headerRequestId, logger.GetRequestId(ctx))
	}
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var requestsDurationSecondsSum = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_request_duration_seconds_sum",
	Help: "Sum of seconds spent on all requests",
})

var requestsDurationSecondsCount = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_request_duration_seconds_count",
	Help: "Count of  all requests",
})

var requestsDurationSecondsBucket = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "http_request_duration_seconds_bucket",
	Help: "group request by tiem repsonses tags",
},
	[]string{"le"})

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		c.Next()

		totalRequests.WithLabelValues(path).Inc()

		time := timer.ObserveDuration()
		incrementRequestOfDuration(time.Seconds())

	}
}
func incrementRequestOfDuration(d float64) {
	go func() {
		requestsDurationSecondsCount.Inc()
		requestsDurationSecondsSum.Add(d)
		if d <= 10 {
			requestsDurationSecondsBucket.WithLabelValues("10").Inc()
		}
		if d <= 1 {
			requestsDurationSecondsBucket.WithLabelValues("1").Inc()
		}
		if d <= 0.3 {
			requestsDurationSecondsBucket.WithLabelValues("0.3").Inc()
		}
		if d <= 0.03 {
			requestsDurationSecondsBucket.WithLabelValues("0.03").Inc()
		}
	}()
}
