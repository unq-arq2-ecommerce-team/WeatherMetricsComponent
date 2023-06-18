package cache

import (
	"context"
	"time"
)

type MemoryCacheClient interface {
	Save(ctx context.Context, table, key string, value interface{}, expiresIn time.Duration) error
	Get(ctx context.Context, table, key string) (interface{}, bool, error)
}
