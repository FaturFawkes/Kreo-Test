package valueobject

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

// Token represents a JWT token
type Token struct {
	value     string
	userID    string
	issuedAt  time.Time
	expiresAt time.Time
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// NewToken creates a new Token value object from a JWT string
func NewToken(tokenString string) (Token, error) {
	normalized := strings.TrimSpace(tokenString)

	// Validate
	if normalized == "" {
		return Token{}, fmt.Errorf("%w: token cannot be empty", ErrInvalidToken)
	}

	parts := strings.Split(normalized, ".")
	if len(parts) != 3 {
		return Token{}, fmt.Errorf("%w: token must have 3 parts", ErrInvalidToken)
	}

	return Token{value: normalized}, nil
}

func NewTokenFromClaims(tokenString string, claims *Claims) Token {
	return Token{
		value:     tokenString,
		userID:    claims.UserID,
		issuedAt:  claims.IssuedAt.Time,
		expiresAt: claims.ExpiresAt.Time,
	}
}

// String returns the token string
func (t Token) String() string {
	return t.value
}

// UserID returns the user ID from the token
func (t Token) UserID() string {
	return t.userID
}

// IssuedAt returns when the token was issued
func (t Token) IssuedAt() time.Time {
	return t.issuedAt
}

// ExpiresAt returns when the token expires
func (t Token) ExpiresAt() time.Time {
	return t.expiresAt
}

// IsExpired checks if the token has expired
func (t Token) IsExpired() bool {
	if t.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(t.expiresAt)
}

// Equals checks if two tokens are equal
func (t Token) Equals(other Token) bool {
	return t.value == other.value
}

// IsEmpty checks if the token is empty
func (t Token) IsEmpty() bool {
	return t.value == ""
}
