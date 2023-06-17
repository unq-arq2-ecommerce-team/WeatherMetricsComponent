package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/logger"
)

const (
	ServiceName = "WeatherMetricsComponent"
)

type Config struct {
	Environment    string               `required:"true" default:"development"`
	Port           int                  `required:"true" default:"8080"`
	PrometheusPort int                  `required:"true" default:"8081"`
	LogLevel       string               `split_words:"true" default:"DEBUG"`
	LokiHost       string               `split_words:"true" required:"true"`
	Redis          RedisConfig          `split_words:"true" required:"true"`
	LocalCache     LocalCacheConfig     `split_words:"true" required:"true"`
	Weather        WeatherEndpoint      `required:"true"`
	CircuitBreaker CircuitBreakerConfig `split_words:"true" required:"true"`
}

type CircuitBreakerConfig struct {
	FailuresRatio float64       `split_words:"true" required:"true"`
	MinRequests   int           `split_words:"true" required:"true"`
	Timeout       time.Duration `split_words:"true" default:"1m"`
}

type WeatherEndpoint struct {
	CurrentTemp EndpointConfig `split_words:"true" required:"true"`
	AvgTemp     EndpointConfig `split_words:"true" required:"true"`
	HttpConfig  HttpConfig     `split_words:"true"`
}

type EndpointConfig struct {
	Url string `split_words:"true" required:"true"`
}

type HttpConfig struct {
	Timeout   time.Duration `default:"10s"`
	Retries   int           `default:"0"`
	RetryWait time.Duration `split_words:"true" default:"15s"`
}

// LocalCacheConfig PurgesExpiration is how often local cache is cleaned up
type LocalCacheConfig struct {
	DefaultExpiration time.Duration `split_words:"true" required:"true"`
	PurgesExpiration  time.Duration `split_words:"true" required:"true"`
}

type RedisConfig struct {
	Uri     string        `split_words:"true" required:"true"`
	Timeout time.Duration `split_words:"true" default:"10s"`
}

func LoadConfig() Config {
	defaultLogger := loggerPkg.DefaultLogger(ServiceName, loggerPkg.JsonFormat)
	// Auto load ".env" file
	err := godotenv.Load()
	if err != nil {
		defaultLogger.Error("error loading .env file")
	}
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err.Error())
	}
	return config
}
