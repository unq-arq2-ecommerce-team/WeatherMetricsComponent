package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	"time"
)

type redisCacheClient struct {
	logger domain.Logger
	client *redis.Client
}

func NewCacheClient(logger domain.Logger, conf config.RedisConfig) cache.MemoryCacheClient {
	return &redisCacheClient{
		logger: logger.WithFields(domain.LoggerFields{"cache.redis": "redisCacheClient"}),
		client: newClient(logger, conf),
	}
}

func newClient(logger domain.Logger, conf config.RedisConfig) *redis.Client {
	ctx, cancelFn := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancelFn()
	opts, err := redis.ParseURL(conf.Uri)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Fatal("error parsing redis url")
	}
	client := redis.NewClient(opts)
	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Fatal("error connecting to redis")
	}

	if conf.OtelEnabled {
		instrumentOtel(client)
		logger.Infof("redis otel instrumentation enabled")
	}

	logger.Infof("connected to redis")
	return client
}

func instrumentOtel(client *redis.Client) {
	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(client); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(client); err != nil {
		panic(err)
	}
}

// Save : returns error if key value cannot be saved in param table
func (c *redisCacheClient) Save(ctx context.Context, table, key string, value interface{}, expiresIn time.Duration) error {
	finalKey := cache.BuildCacheKey(table, key)
	log := c.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"method": "Save", "table": table, "key": key, "finalKey": finalKey})

	value, err := json.Marshal(value)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("cannot marshal value")
		return err
	}

	if err := c.client.Set(ctx, finalKey, value, expiresIn).Err(); err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error save in redis with key %s", finalKey)
		return err
	}
	log.Debugf("redis saved ok with final key %s", finalKey)
	return nil
}

// Get : returns element if was found otherwise nil; a boolean if key was found; and error if something went wrong
func (c *redisCacheClient) Get(ctx context.Context, table, key string) (interface{}, bool, error) {
	finalKey := cache.BuildCacheKey(table, key)
	log := c.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"method": "Get", "table": table, "key": key, "finalKey": finalKey})
	elemString, err := c.client.Get(ctx, finalKey).Result()
	if err != nil && err != redis.Nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error get in redis with key %s", finalKey)
		return nil, false, err
	}
	found := err != redis.Nil
	log.Debugf("redis get ok with final key %s, found %v and value %s", finalKey, elemString)
	return elemString, found, nil
}
