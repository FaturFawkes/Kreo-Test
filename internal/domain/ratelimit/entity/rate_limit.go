package entity

import (
	"time"
	"upwork-test/internal/domain/ratelimit/valueobject"
)

// RateLimit represents a rate limit record for a user
type RateLimit struct {
	UserID       string
	Tier         valueobject.RateLimitTier
	WindowSize   time.Duration
	MaxRequests  int
	RequestCount int
	WindowStart  time.Time
}

// NewRateLimit creates a new RateLimit entity
func NewRateLimit(userID string, tier valueobject.RateLimitTier) *RateLimit {
	maxRequests := tier.MaxRequests()
	windowSize := tier.Window()

	return &RateLimit{
		UserID:       userID,
		Tier:         tier,
		WindowSize:   windowSize,
		MaxRequests:  maxRequests,
		RequestCount: 0,
		WindowStart:  time.Now(),
	}
}

// IsAllowed checks if a new request is allowed under the rate limit
func (rl *RateLimit) IsAllowed() bool {
	return rl.RequestCount < rl.MaxRequests
}

// IncrementCount increments the request count
func (rl *RateLimit) IncrementCount() {
	rl.RequestCount++
}

// RemainingRequests returns the number of requests remaining in the current window
func (rl *RateLimit) RemainingRequests() int {
	remaining := rl.MaxRequests - rl.RequestCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ResetTime returns when the rate limit window resets
func (rl *RateLimit) ResetTime() time.Time {
	return rl.WindowStart.Add(rl.WindowSize)
}

// IsExpired checks if the current window has expired
func (rl *RateLimit) IsExpired() bool {
	return time.Now().After(rl.ResetTime())
}

// Reset resets the rate limit counter for a new window
func (rl *RateLimit) Reset() {
	rl.RequestCount = 0
	rl.WindowStart = time.Now()
}
