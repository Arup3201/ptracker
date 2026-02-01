package infra

import (
	"context"
	"time"

	"github.com/ptracker/internal/interfaces"
	"github.com/redis/go-redis/v9"
)

type inMemory struct {
	client *redis.Client
}

func NewInMemory(client *redis.Client) interfaces.InMemory {
	return &inMemory{client: client}
}

func (c *inMemory) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *inMemory) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *inMemory) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
