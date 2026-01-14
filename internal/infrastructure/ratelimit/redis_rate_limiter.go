package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"upwork-test/internal/domain/ratelimit/entity"
	"upwork-test/internal/domain/ratelimit/valueobject"
	"upwork-test/internal/infrastructure/cache"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter implements the rate limiting using Redis sliding window.
type RedisRateLimiter struct {
	client     *redis.Client
	keyBuilder *cache.KeyBuilder
}

// NewRedisRateLimiter creates a new Redis-backed rate limiter.
func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:     client,
		keyBuilder: cache.NewKeyBuilder("kalshi"),
	}
}

// Get retrieves the current rate limit state for a user.
func (r *RedisRateLimiter) Get(ctx context.Context, userID string) (*entity.RateLimit, error) {
	key := r.keyBuilder.RateLimitCounter(userID, "current")

	countStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit: %w", err)
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, fmt.Errorf("invalid count value: %w", err)
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get TTL: %w", err)
	}

	windowStart := time.Now().Add(-time.Minute).Add(ttl)

	// We can't fully reconstruct the RateLimit without tier information
	// This is a limitation of the Get method - in practice, the tier should be passed
	// For now, we'll return a basic structure
	return &entity.RateLimit{
		UserID:       userID,
		RequestCount: count,
		WindowStart:  windowStart,
	}, nil
}

// Save persists the rate limit state.
func (r *RedisRateLimiter) Save(ctx context.Context, rateLimit *entity.RateLimit) error {
	key := r.keyBuilder.RateLimitCounter(rateLimit.UserID, "current")

	ttl := time.Until(rateLimit.ResetTime())
	if ttl < 0 {
		ttl = rateLimit.WindowSize
	}

	err := r.client.Set(ctx, key, rateLimit.RequestCount, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save rate limit: %w", err)
	}

	return nil
}

// IncrementAndCheck atomically increments and checks if request is allowed.
// This uses a Lua script to ensure atomicity.
func (r *RedisRateLimiter) IncrementAndCheck(ctx context.Context, userID string, tier valueobject.RateLimitTier) (bool, int, time.Time, error) {
	// Use tier name as the window identifier to separate different tiers
	key := r.keyBuilder.RateLimitCounter(userID, tier.Name())
	maxRequests := tier.MaxRequests()
	windowSeconds := int(tier.Window().Seconds())

	script := redis.NewScript(`
		local key = KEYS[1]
		local max_requests = tonumber(ARGV[1])
		local window_seconds = tonumber(ARGV[2])
		
		local count = redis.call('INCR', key)
		local ttl = redis.call('TTL', key)
		
		if ttl == -1 then
			redis.call('EXPIRE', key, window_seconds)
			ttl = window_seconds
		end
		
		return {count, ttl}
	`)

	result, err := script.Run(ctx, r.client, []string{key}, maxRequests, windowSeconds).Result()
	if err != nil {
		return false, 0, time.Time{}, fmt.Errorf("failed to increment rate limit: %w", err)
	}

	resultSlice := result.([]interface{})
	count := int(resultSlice[0].(int64))
	ttl := int(resultSlice[1].(int64))

	allowed := count <= maxRequests
	remaining := maxRequests - count
	if remaining < 0 {
		remaining = 0
	}

	resetTime := time.Now().Add(time.Duration(ttl) * time.Second)

	return allowed, remaining, resetTime, nil
}
