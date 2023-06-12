package cache

import (
	"context"
	"fmt"
	c "github.com/patrickmn/go-cache"
	"github.com/unq-arq2-ecommerce-team/WeatherMetricsComponent/internal/domain"
	"time"
)

type MemoryCacheClient interface {
	Save(ctx context.Context, table, key string, value interface{}, expiresIn time.Duration) error
	Get(ctx context.Context, table, key string) (interface{}, bool, error)
}

type localMemoryCacheClient struct {
	logger      domain.Logger
	cacheTables map[string]*c.Cache
}

// NewLocalMemoryCacheClient : must be used only for services with one host
func NewLocalMemoryCacheClient(logger domain.Logger, cacheTables map[string]*c.Cache) MemoryCacheClient {
	return &localMemoryCacheClient{
		logger:      logger.WithFields(domain.LoggerFields{"cache": "localMemoryCacheClient"}),
		cacheTables: cacheTables,
	}
}

// Save : returns error if key value cannot be saved in param table
func (c *localMemoryCacheClient) Save(ctx context.Context, table, key string, value interface{}, expiresIn time.Duration) error {
	log := c.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"method": "Save", "table": table, "key": key})
	cacheTable, err := c.getCacheTable(log, table)
	if err != nil {
		return err
	}
	cacheTable.Set(key, value, expiresIn)
	log.Debugf("cache save ok with table %s and key %s", table, key)
	return nil
}

// Get : returns element if was found otherwise nil; a boolean if key was found; and error if something went wrong
func (c *localMemoryCacheClient) Get(ctx context.Context, table, key string) (interface{}, bool, error) {
	log := c.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"method": "Get", "table": table, "key": key})
	cacheTable, err := c.getCacheTable(log, table)
	if err != nil {
		return nil, false, err
	}
	elem, found := cacheTable.Get(key)
	log.Debugf("cache get ok with table %s and key %s", table, key)
	return elem, found, nil
}

func (c *localMemoryCacheClient) getCacheTable(log domain.Logger, table string) (*c.Cache, error) {
	cacheTable, found := c.cacheTables[table]
	if !found {
		log.Errorf("cache table %s not exist", table)
		return nil, fmt.Errorf("cache table %s not exist", table)
	}
	return cacheTable, nil
}
