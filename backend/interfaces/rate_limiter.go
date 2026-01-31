package interfaces

import "context"

type RateLimiter interface {
	Allow(ctx context.Context, userId string) (bool, error)
	Tokens(ctx context.Context,
		key string) (int, error)
	Capacity(ctx context.Context) int
	RetryAfter(ctx context.Context, userId string) (int, error)
}
