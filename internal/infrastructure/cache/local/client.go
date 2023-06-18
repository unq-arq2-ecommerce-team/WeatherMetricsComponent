package local

import (
	"context"
	c "github.com/patrickmn/go-cache"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/cache"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/infrastructure/config"
	"time"
)

type localMemoryCacheClient struct {
	logger        domain.Logger
	cacheInstance *c.Cache
}

// NewLocalMemoryCacheClient : must be used only for services with one host
func NewLocalMemoryCacheClient(logger domain.Logger, conf config.LocalCacheConfig) cache.MemoryCacheClient {
	return &localMemoryCacheClient{
		logger:        logger.WithFields(domain.LoggerFields{"cache.local": "localMemoryCacheClient"}),
		cacheInstance: c.New(conf.DefaultExpiration, conf.PurgesExpiration),
	}
}

// Save : returns error if key value cannot be saved in param table
func (c *localMemoryCacheClient) Save(ctx context.Context, table, key string, value interface{}, expiresIn time.Duration) error {
	finalKey := cache.BuildCacheKey(table, key)
	log := c.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"method": "Save", "table": table, "key": key, "finalKey": finalKey})
	c.cacheInstance.Set(finalKey, value, expiresIn)
	log.Debugf("local cache saved ok with final key %s", finalKey)
	return nil
}

// Get : returns element if was found otherwise nil; a boolean if key was found; and error if something went wrong
func (c *localMemoryCacheClient) Get(ctx context.Context, table, key string) (interface{}, bool, error) {
	finalKey := cache.BuildCacheKey(table, key)
	log := c.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"method": "Get", "table": table, "key": key, "finalKey": finalKey})
	elem, found := c.cacheInstance.Get(finalKey)
	log.Debugf("cache get ok with final key %s", finalKey)
	return elem, found, nil
}
