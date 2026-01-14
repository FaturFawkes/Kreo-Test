package usecase

import (
	"errors"
	"fmt"
	"upwork-test/internal/application/dto"
	"upwork-test/internal/domain/auth/entity"
	"upwork-test/internal/domain/auth/service"
	"upwork-test/internal/domain/auth/valueobject"
)

var (
	// ErrInvalidAPIKey is returned when API key is invalid
	ErrInvalidAPIKey = errors.New("invalid API key")
)

// Authenticate use case handles user authentication
type Authenticate struct {
	tokenService *service.TokenService
	validAPIKey string
}

// NewAuthenticate creates a new Authenticate use case
func NewAuthenticate(tokenService *service.TokenService, validAPIKey string) *Authenticate {
	return &Authenticate{
		tokenService: tokenService,
		validAPIKey:  validAPIKey,
	}
}

// Execute authenticates a user and returns a JWT token
func (uc *Authenticate) Execute(request *dto.AuthRequest) (*dto.TokenResponse, error) {
	credentials, err := valueobject.NewCredentials(request.APIKey)
	if err != nil {
		return nil, ErrInvalidAPIKey
	}

	if credentials.APIKey() != uc.validAPIKey {
		return nil, ErrInvalidAPIKey
	}

	user := entity.NewUser("default-user", credentials)
	user.UpdateLastLogin()

	token, err := uc.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return dto.NewTokenResponse(token.String(), token.ExpiresAt()), nil
}
