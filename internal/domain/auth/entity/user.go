package entity

import (
	"time"
	"upwork-test/internal/domain/auth/valueobject"
)

// User represents an authenticated user
type User struct {
	ID          string
	Credentials valueobject.Credentials
	CreatedAt   time.Time
	LastLogin   time.Time
}

// NewUser creates a new User entity
func NewUser(id string, credentials valueobject.Credentials) *User {
	now := time.Now()
	return &User{
		ID:          id,
		Credentials: credentials,
		CreatedAt:   now,
		LastLogin:   now,
	}
}

// UpdateLastLogin updates the user's last login time
func (u *User) UpdateLastLogin() {
	u.LastLogin = time.Now()
}

// IsCredentialsValid checks if the provided API key matches
func (u *User) IsCredentialsValid(apiKey string) bool {
	return u.Credentials.APIKey() == apiKey
}

