package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"upwork-test/internal/domain/auth/service"
	ratelimit "upwork-test/internal/domain/ratelimit/service"
	"upwork-test/internal/domain/ratelimit/valueobject"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware creates a middleware that enforces rate limits.
func RateLimitMiddleware(limiter *ratelimit.RateLimiter, tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine user ID and tier
		userID, tier := getUserIDAndTier(c, tokenService)

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

// getUserIDAndTier extracts user ID and rate limit tier from the request.
// It attempts to decode the JWT token from the Authorization header to determine
// if the user is authenticated. This runs before the Auth middleware, so it doesn't
// rely on context values.
func getUserIDAndTier(c *gin.Context, tokenService *service.TokenService) (string, valueobject.RateLimitTier) {
	// Try to extract and validate JWT token
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]

			// Attempt to validate token
			token, err := tokenService.ValidateToken(tokenString)
			if err == nil {
				// Valid token - use authenticated tier
				userID := token.UserID()
				c.Set("user_id", userID) // Pre-set for Auth middleware
				return userID, valueobject.Authenticated
			}
			// Invalid/expired token - will be caught by Auth middleware later
			// For rate limiting purposes, treat as unauthenticated
		}
	}

	// For unauthenticated users, use client IP
	clientIP := c.ClientIP()
	return clientIP, valueobject.Unauthenticated
}
