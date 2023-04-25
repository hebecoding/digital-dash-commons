package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache interface {
	Get(ctx context.Context, key string) (string, error)
	GetDelete(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type Cache struct {
	cache *redis.Client
}

func NewCache(_ context.Context, cache *redis.Client) *Cache {
	return &Cache{
		cache: cache,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.cache.Get(ctx, key).Result()
}

func (c *Cache) GetDelete(ctx context.Context, key string) (string, error) {
	return c.cache.GetDel(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value string, duration time.Duration) error {
	return c.cache.Set(ctx, key, value, duration).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.cache.Del(ctx, key).Err()
}
