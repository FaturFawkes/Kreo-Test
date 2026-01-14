package response

import "time"

// TokenResponse represents the JWT token response
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	TokenType string    `json:"token_type"`
}

// NewTokenResponse creates a new token response
func NewTokenResponse(token string, expiresAt time.Time) *TokenResponse {
	return &TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		TokenType: "Bearer",
	}
}
