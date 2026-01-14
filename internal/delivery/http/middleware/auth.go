package middleware

import (
	"net/http"
	"strings"

	"upwork-test/internal/delivery/http/response"
	"upwork-test/internal/domain/auth/service"
	"upwork-test/internal/domain/auth/valueobject"

	"github.com/gin-gonic/gin"
)

// Auth returns a middleware that validates JWT tokens
func Auth(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID, _ := c.Get("trace_id")

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse(
				http.StatusUnauthorized,
				"Missing Authorization header",
				traceID.(string),
			))
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse(
				http.StatusUnauthorized,
				"Invalid Authorization header format",
				traceID.(string),
			))
			return
		}

		tokenString := parts[1]

		// Validate token
		token, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			if err == valueobject.ErrTokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse(
					http.StatusUnauthorized,
					"Token expired",
					traceID.(string),
				))
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse(
				http.StatusUnauthorized,
				"Invalid token",
				traceID.(string),
			))
			return
		}

		// Store user ID in context
		c.Set("user_id", token.UserID())

		c.Next()
	}
}
