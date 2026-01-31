package db

import (
	"context"
	"time"

	"github.com/ptracker/interfaces"
	"github.com/redis/go-redis/v9"
)

type redisInMemory struct {
	client *redis.Client
}

func NewRedisInMemory(client *redis.Client) interfaces.InMemory {
	return &redisInMemory{client: client}
}

func (c *redisInMemory) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisInMemory) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *redisInMemory) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
