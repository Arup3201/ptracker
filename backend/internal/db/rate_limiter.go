package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ptracker/internal/interfaces"
	"github.com/redis/go-redis/v9"
)

type redisRateLimiter struct {
	client         *redis.Client
	capacity, rate int
	retry          int64
	tbFunc         *redis.Script
}

func NewRedisRateLimiter(client *redis.Client, cap, rate int) interfaces.RateLimiter {
	function := redis.NewScript(RedisLuaScript)

	return &redisRateLimiter{
		client:   client,
		capacity: cap,
		rate:     rate,
		tbFunc:   function,
	}
}

func (rl *redisRateLimiter) Allow(ctx context.Context,
	key string) (bool, error) {
	value, err := rl.tbFunc.Run(
		ctx,
		rl.client,
		[]string{key},
		rl.Capacity,
		rl.rate,
	).Result()
	if err != nil {
		return false, fmt.Errorf("allow request run lua script: %w", err)
	}

	if value.(int) == 1 {
		return true, nil
	}

	return false, nil
}

func (rl *redisRateLimiter) Tokens(ctx context.Context,
	key string) (int, error) {

	value, err := rl.client.HGet(ctx, key, "tokens").Result()
	if err != nil {
		return 0, fmt.Errorf("client hget: %w", err)
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("string parse error: %w", err)
	}

	return v, nil
}

func (rl *redisRateLimiter) Capacity(ctx context.Context) int {
	return rl.capacity
}

func (rl *redisRateLimiter) RetryAfter(ctx context.Context,
	key string) (int, error) {

	value, err := rl.client.HGet(ctx, key, "retry_after").Result()
	if err != nil {
		return 0, fmt.Errorf("client hget: %w", err)
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("string parse error: %w", err)
	}

	return v, nil
}

const RedisLuaScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])

local now = redis.call("TIME")
local now_sec = tonumber(now[1])
local now_usec = tonumber(now[2])
local now_ts = now_sec + now_usec / 1e6

local tokens = tonumber(redis.call("HGET", key, "tokens"))
if tokens == nil then
	tokens = capacity -- start full
end

local last_refill = tonumber(redis.call("HGET", key, "last_refill"))
if last_refill == nil then
	last_refill = now_ts
end

-- refill
local elapsed = now_ts - last_refill
local tokens_to_add = math.floor(elapsed) * rate

if tokens_to_add > 0 then
	tokens = math.min(capacity, tokens+tokens_to_add)
	redis.call("HSET", key, "last_refill", now_ts)
end

-- deny if not enough tokens
if tokens < 1 then
    local retry_after = math.ceil(((1 - (elapsed*rate)) / rate)*1000)
	redis.call("HSET", key, "retry_after", retry_after)
    return 0
end

-- consume
tokens = tokens - 1
redis.call("HSET", key, "tokens", tokens)

return 1
`
