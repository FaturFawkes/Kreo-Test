package usecase

import (
	"context"
	"errors"
	"fmt"
	"upwork-test/internal/application/dto"
	"upwork-test/internal/domain/category/entity"
	"upwork-test/internal/domain/category/repository"
)

var (
	// ErrCategoryOverviewNotFound is returned when category overview is not found
	ErrCategoryOverviewNotFound = errors.New("category overview not found")
)

// GetCategoryOverview use case retrieves overview metrics for a category.
type GetCategoryOverview struct {
	categoryRepo repository.CategoryRepository
}

// NewGetCategoryOverview creates a new GetCategoryOverview use case.
func NewGetCategoryOverview(categoryRepo repository.CategoryRepository) *GetCategoryOverview {
	return &GetCategoryOverview{
		categoryRepo: categoryRepo,
	}
}

// Execute retrieves category overview metrics.
func (uc *GetCategoryOverview) Execute(ctx context.Context, categoryName string) (*dto.CategoryOverviewDTO, error) {
	if categoryName == "" {
		return nil, errors.New("category name cannot be empty")
	}

	overview, err := uc.categoryRepo.GetOverview(ctx, categoryName)
	if err != nil {
		if errors.Is(err, repository.ErrOverviewNotFound) {
			return nil, ErrCategoryOverviewNotFound
		}
		return nil, fmt.Errorf("failed to get category overview: %w", err)
	}

	return uc.overviewToDTO(overview), nil
}

// overviewToDTO converts a CategoryOverview entity to DTO.
func (uc *GetCategoryOverview) overviewToDTO(overview *entity.CategoryOverview) *dto.CategoryOverviewDTO {
	return &dto.CategoryOverviewDTO{
		CategoryName:     overview.CategoryName.String(),
		TotalMarkets:     overview.TotalMarkets,
		TotalVolume24h:   overview.TotalVolume24h,
		AverageLiquidity: overview.AverageLiquidity,
		ActiveTraders24h: overview.ActiveTraders24h,
		ComputedAt:       overview.ComputedAt,
		ExpiresAt:        overview.ExpiresAt,
	}
}
