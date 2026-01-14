package service

import (
	"fmt"
	"time"
	"upwork-test/internal/domain/auth/valueobject"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService handles JWT token generation and validation
type TokenService struct {
	secretKey  string
	expiration time.Duration
}

// NewTokenService creates a new TokenService
func NewTokenService(secretKey string, expiration time.Duration) *TokenService {
	return &TokenService{
		secretKey:  secretKey,
		expiration: expiration,
	}
}

// GenerateToken generates a new JWT token for a user
func (s *TokenService) GenerateToken(userID string) (valueobject.Token, error) {
	now := time.Now()
	expiresAt := now.Add(s.expiration)

	claims := &valueobject.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return valueobject.Token{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return valueobject.NewTokenFromClaims(tokenString, claims), nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *TokenService) ValidateToken(tokenString string) (valueobject.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &valueobject.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return valueobject.Token{}, fmt.Errorf("%w: %v", valueobject.ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*valueobject.Claims)
	if !ok || !token.Valid {
		return valueobject.Token{}, valueobject.ErrInvalidToken
	}

	return valueobject.NewTokenFromClaims(tokenString, claims), nil
}

// ExtractUserID extracts the user ID from a token without full validation
// Useful for logging or metrics where we need the user ID even if token is expired
func (s *TokenService) ExtractUserID(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &valueobject.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*valueobject.Claims)
	if !ok {
		return "", valueobject.ErrInvalidToken
	}

	return claims.UserID, nil
}