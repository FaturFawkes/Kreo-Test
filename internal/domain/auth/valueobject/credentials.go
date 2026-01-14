package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Credentials struct {
	apiKey string
}

// NewCredentials creates a new Credentials value object
func NewCredentials(apiKey string) (Credentials, error) {
	normalized := strings.TrimSpace(apiKey)

	// Validate
	if normalized == "" {
		return Credentials{}, fmt.Errorf("%w: API key cannot be empty", ErrInvalidCredentials)
	}

	if len(normalized) > 256 {
		return Credentials{}, fmt.Errorf("%w: API key too long (maximum 256 characters)", ErrInvalidCredentials)
	}

	return Credentials{apiKey: normalized}, nil
}

// APIKey returns the API key value
func (c Credentials) APIKey() string {
	return c.apiKey
}

// Equals checks if two credentials are equal
func (c Credentials) Equals(other Credentials) bool {
	return c.apiKey == other.apiKey
}

// IsEmpty checks if the credentials are empty
func (c Credentials) IsEmpty() bool {
	return c.apiKey == ""
}

// MaskedAPIKey returns a masked version of the API key for logging
func (c Credentials) MaskedAPIKey() string {
	if len(c.apiKey) <= 8 {
		return "****"
	}
	return c.apiKey[:4] + "****" + c.apiKey[len(c.apiKey)-4:]
}
