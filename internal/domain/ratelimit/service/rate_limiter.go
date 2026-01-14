package service

import (
	"context"
	"time"

	"upwork-test/internal/domain/ratelimit/entity"
	"upwork-test/internal/domain/ratelimit/valueobject"
)

// RateLimitRepository defines the interface for rate limit persistence.
type RateLimitRepository interface {
	// Get retrieves the current rate limit state for a user
	Get(ctx context.Context, userID string) (*entity.RateLimit, error)

	// Save persists the rate limit state
	Save(ctx context.Context, rateLimit *entity.RateLimit) error

	// IncrementAndCheck atomically increments and checks if request is allowed
	IncrementAndCheck(ctx context.Context, userID string, tier valueobject.RateLimitTier) (allowed bool, remaining int, resetTime time.Time, err error)
}

// RateLimiter is a domain service that handles rate limiting logic.
type RateLimiter struct {
	repo RateLimitRepository
}

// NewRateLimiter creates a new RateLimiter service.
func NewRateLimiter(repo RateLimitRepository) *RateLimiter {
	return &RateLimiter{
		repo: repo,
	}
}

// CheckLimit verifies if a request is allowed for the user and increments the counter.
// Returns whether the request is allowed, remaining requests, and reset time.
func (rl *RateLimiter) CheckLimit(ctx context.Context, userID string, tier valueobject.RateLimitTier) (allowed bool, remaining int, resetTime time.Time, err error) {
	return rl.repo.IncrementAndCheck(ctx, userID, tier)
}

// GetCurrentLimit retrieves the current rate limit status for a user without incrementing.
func (rl *RateLimiter) GetCurrentLimit(ctx context.Context, userID string) (*entity.RateLimit, error) {
	return rl.repo.Get(ctx, userID)
}
