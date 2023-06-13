package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
	"time"
)

const ServiceName = "WeatherMetricsComponent"

type Config struct {
	Environment    string          `required:"true" default:"development"`
	Port           int             `required:"true" default:"8080"`
	PrometheusPort int             `required:"true" default:"8081"`
	LogLevel       string          `split_words:"true" default:"DEBUG"`
	Weather        WeatherEndpoint `required:"true"`
}

type WeatherEndpoint struct {
	CurrentTemp EndpointConfig `split_words:"true" required:"true"`
	AvgTemp     EndpointConfig `split_words:"true" required:"true"`
	HttpConfig  HttpConfig     `split_words:"true"`
}

type EndpointConfig struct {
	Url   string      `split_words:"true" required:"true"`
	Cache CacheConfig `split_words:"true" required:"true"`
}

type HttpConfig struct {
	Timeout   time.Duration `default:"10s"`
	Retries   int           `default:"0"`
	RetryWait time.Duration `split_words:"true" default:"15s"`
}

type CacheConfig struct {
	DefaultExp time.Duration `split_words:"true" required:"true"`
	PurgesExp  time.Duration `split_words:"true" required:"true"`
}

func LoadConfig() Config {
	primitiveLogger := loggerPkg.New(&loggerPkg.Config{
		ServiceName: ServiceName,
		LogFormat:   loggerPkg.JsonFormat,
	})

	// Auto load ".env" file
	err := godotenv.Load()
	if err != nil {
		primitiveLogger.Error("error loading .env file")
	}
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err.Error())
	}
	return config
}
