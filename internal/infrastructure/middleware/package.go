package middleware

import (
	"github.com/gin-gonic/gin"
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
