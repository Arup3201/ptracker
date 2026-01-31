package services

import (
	"context"

	"github.com/ptracker/interfaces"
)

func getBucketKey(userId string) string {
	return "bucket:user:" + userId
}

type limiterService struct {
	store interfaces.Store
}

func NewLimiterService(store interfaces.Store) interfaces.LimiterService {
	return &limiterService{
		store: store,
	}
}

func (s *limiterService) IsAllowed(ctx context.Context,
	userId string) (bool, error) {

	key := getBucketKey(userId)
	return s.store.RateLimiter().Allow(ctx, key)
}

func (s *limiterService) GetTokens(ctx context.Context,
	userId string) (int, error) {

	key := getBucketKey(userId)
	return s.store.RateLimiter().Tokens(ctx, key)
}

func (s *limiterService) GetCapacity(ctx context.Context) int {

	return s.store.RateLimiter().Capacity(ctx)
}

func (s *limiterService) GetRetryTime(ctx context.Context,
	userId string) (int, error) {

	key := getBucketKey(userId)
	return s.store.RateLimiter().RetryAfter(ctx, key)
}
