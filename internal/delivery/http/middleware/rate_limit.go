package middleware

import (
	"net/http"
	"strconv"
	"upwork-test/internal/domain/ratelimit/service"
	"upwork-test/internal/domain/ratelimit/valueobject"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware creates a middleware that enforces rate limits.
func RateLimitMiddleware(limiter *service.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine user ID and tier
		userID, tier := getUserIDAndTier(c)

		// Check rate limit
		allowed, remaining, resetTime, err := limiter.CheckLimit(c.Request.Context(), userID, tier)
		if err != nil {
			// Log error but don't block the request on rate limit check failure
			// This ensures availability over strict rate limiting
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(tier.MaxRequests()))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		// If rate limit exceeded, return 429
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate_limit_exceeded",
				"message":     "Rate limit exceeded. Please try again later.",
				"retry_after": resetTime.Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getUserIDAndTier extracts user ID and rate limit tier from the request context.
// If user is authenticated, use their user ID with authenticated tier.
// Otherwise, use IP address with unauthenticated tier.
func getUserIDAndTier(c *gin.Context) (string, valueobject.RateLimitTier) {
	// Check if user is authenticated (set by Auth middleware)
	userID, exists := c.Get("user_id")
	if exists {
		// For authenticated users
		return userID.(string), valueobject.Authenticated
	}

	// For unauthenticated users, use client IP
	clientIP := c.ClientIP()
	return clientIP, valueobject.Unauthenticated
}
