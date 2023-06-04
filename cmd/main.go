package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	infra "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure"
)

func main() {
	logger := log.Default()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}
	app := infra.NewGinApplication()
	logger.Fatal(app.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
