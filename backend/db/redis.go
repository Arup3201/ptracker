package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInMemory struct {
	client *redis.Client
}

func NewRedisInMemory(client *redis.Client) *RedisInMemory {
	return &RedisInMemory{client: client}
}

func (c *RedisInMemory) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *RedisInMemory) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisInMemory) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
