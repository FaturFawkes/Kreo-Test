package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logging returns a middleware that logs HTTP requests
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate trace ID if not present
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set("trace_id", traceID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Build log fields
		fields := map[string]interface{}{
			"trace_id":    traceID,
			"method":      c.Request.Method,
			"path":        path,
			"query":       query,
			"status_code": statusCode,
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}

		// Add error if present
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// Log based on status code
		// Note: logmanager uses context-based logging, not direct Info/Warn/Error
		// Using fmt.Printf for now as a placeholder
		if statusCode >= 500 {
			fmt.Printf("[ERROR] HTTP request completed with server error: %v\n", fields)
		} else if statusCode >= 400 {
			fmt.Printf("[WARN] HTTP request completed with client error: %v\n", fields)
		} else {
			fmt.Printf("[INFO] HTTP request completed: %v\n", fields)
		}
	}
}
