package valueobject

import (
	"fmt"
	"time"
)

// RateLimitTier represents different rate limit tiers in the system.
type RateLimitTier struct {
	name        string
	maxRequests int
	window      time.Duration
}

// Predefined rate limit tiers.
var (
	// Authenticated tier: 100 requests per minute for logged-in users
	Authenticated = RateLimitTier{
		name:        "authenticated",
		maxRequests: 100,
		window:      time.Minute,
	}

	// Unauthenticated tier: 10 requests per minute for anonymous users
	Unauthenticated = RateLimitTier{
		name:        "unauthenticated",
		maxRequests: 10,
		window:      time.Minute,
	}

	// Worker tier: 80 requests per minute for background workers
	Worker = RateLimitTier{
		name:        "worker",
		maxRequests: 80,
		window:      time.Minute,
	}
)

// NewRateLimitTier creates a new RateLimitTier from a tier name.
func NewRateLimitTier(tierName string) (RateLimitTier, error) {
	switch tierName {
	case "authenticated":
		return Authenticated, nil
	case "unauthenticated":
		return Unauthenticated, nil
	case "worker":
		return Worker, nil
	default:
		return RateLimitTier{}, fmt.Errorf("invalid rate limit tier: %s", tierName)
	}
}

// Name returns the tier name.
func (t RateLimitTier) Name() string {
	return t.name
}

// MaxRequests returns the maximum number of requests allowed in the window.
func (t RateLimitTier) MaxRequests() int {
	return t.maxRequests
}

// Window returns the time window duration.
func (t RateLimitTier) Window() time.Duration {
	return t.window
}

// String returns a string representation of the tier.
func (t RateLimitTier) String() string {
	return fmt.Sprintf("%s (%d req/%s)", t.name, t.maxRequests, t.window)
}
