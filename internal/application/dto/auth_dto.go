package dto

import "time"

// AuthRequest represents an authentication request
type AuthRequest struct {
	APIKey string `json:"api_key"`
}

// TokenResponse represents a JWT token response
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	TokenType string    `json:"token_type"`
}

// NewTokenResponse creates a new TokenResponse
func NewTokenResponse(token string, expiresAt time.Time) *TokenResponse {
	return &TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		TokenType: "Bearer",
	}
}
