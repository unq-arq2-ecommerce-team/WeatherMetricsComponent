package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginApplication struct {
}

func NewGinApplication() *ginApplication {
	return &ginApplication{}
}

func (ginApplication *ginApplication) Run(addr string) error {

	router := gin.Default()
	router.GET("/health", HealthCheck)
	return router.Run(addr)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckRes{Data: "Server is up and running"})
}

type HealthCheckRes struct {
	Data string `json:"data"`
}
