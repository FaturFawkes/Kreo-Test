package middleware

import (
	"net/http"
	"upwork-test/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler returns a middleware that handles errors consistently
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last()

		// Get trace ID from context
		traceID, _ := c.Get("trace_id")
		traceIDStr := ""
		if traceID != nil {
			traceIDStr = traceID.(string)
		}

		// Determine status code (if not already set)
		statusCode := c.Writer.Status()
		if statusCode == http.StatusOK {
			statusCode = http.StatusInternalServerError
		}

		// Build error response
		errorResponse := response.NewErrorResponse(
			statusCode,
			err.Error(),
			traceIDStr,
		)

		// Send error response
		c.JSON(statusCode, errorResponse)
	}
}
