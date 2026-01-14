package repository

import (
	"context"
	"errors"

	"upwork-test/internal/domain/category/entity"
)

var (
	// ErrCategoryNotFound is returned when a category is not found
	ErrCategoryNotFound = errors.New("category not found")
	// ErrOverviewNotFound is returned when category overview is not found
	ErrOverviewNotFound = errors.New("category overview not found")
)

// CategoryRepository defines the interface for category data access.
type CategoryRepository interface {
	// GetAll retrieves all available categories
	GetAll(ctx context.Context) ([]*entity.Category, error)

	// GetByName retrieves a single category by name
	GetByName(ctx context.Context, name string) (*entity.Category, error)

	// GetOverview retrieves the overview metrics for a category (cached for 10 minutes)
	GetOverview(ctx context.Context, categoryName string) (*entity.CategoryOverview, error)

	// SaveOverview saves or updates category overview metrics
	SaveOverview(ctx context.Context, overview *entity.CategoryOverview) error
}
