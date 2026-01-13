package response

import (
	"net/http"
	"time"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	TraceID   string      `json:"trace_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(statusCode int, message string, traceID string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    statusCode,
			Message: message,
			Type:    getErrorType(statusCode),
		},
		TraceID:   traceID,
		Timestamp: time.Now().Unix(),
	}
}

// getErrorType returns a human-readable error type based on status code
func getErrorType(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "internal_server_error"
	case statusCode == http.StatusNotFound:
		return "not_found"
	case statusCode == http.StatusUnauthorized:
		return "unauthorized"
	case statusCode == http.StatusForbidden:
		return "forbidden"
	case statusCode == http.StatusBadRequest:
		return "bad_request"
	case statusCode == http.StatusTooManyRequests:
		return "rate_limit_exceeded"
	case statusCode >= 400:
		return "client_error"
	default:
		return "unknown_error"
	}
}

// ValidationErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Error     ValidationErrorDetail `json:"error"`
	TraceID   string                `json:"trace_id,omitempty"`
	Timestamp int64                 `json:"timestamp"`
}

// ValidationErrorDetail contains validation error details
type ValidationErrorDetail struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Type    string            `json:"type"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(fieldErrors map[string]string, traceID string) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		Error: ValidationErrorDetail{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Type:    "validation_error",
			Fields:  fieldErrors,
		},
		TraceID:   traceID,
		Timestamp: time.Now().Unix(),
	}
}
