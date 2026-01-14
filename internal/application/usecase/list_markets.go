package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"upwork-test/internal/application/dto"
	"upwork-test/internal/domain/market/entity"
	"upwork-test/internal/domain/market/repository"
)

var (
	// ErrInvalidPage is returned when page number is invalid
	ErrInvalidPage = errors.New("invalid page number")
	// ErrInvalidLimit is returned when limit is invalid
	ErrInvalidLimit = errors.New("invalid limit")
	// ErrCategoryNotFound is returned when category is not found
	ErrCategoryNotFound = errors.New("category not found")
)

// ListMarkets use case retrieves paginated markets for a category.
type ListMarkets struct {
	marketRepo repository.MarketRepository
}

// NewListMarkets creates a new ListMarkets use case.
func NewListMarkets(marketRepo repository.MarketRepository) *ListMarkets {
	return &ListMarkets{
		marketRepo: marketRepo,
	}
}

// Execute retrieves markets for a category with pagination.
func (uc *ListMarkets) Execute(ctx context.Context, category string, page int, limit int, status string) (*dto.MarketListDTO, error) {
	if page < 1 {
		return nil, ErrInvalidPage
	}
	if limit < 1 || limit > 100 {
		return nil, ErrInvalidLimit
	}

	markets, total, err := uc.marketRepo.ListByCategory(ctx, category, page, limit, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list markets: %w", err)
	}

	marketDTOs := make([]*dto.MarketDTO, len(markets))
	for i, market := range markets {
		marketDTOs[i] = uc.marketToDTO(market)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	var nextURL, prevURL string
	if page < totalPages {
		nextURL = fmt.Sprintf("/api/v1/categories/%s/markets?page=%d&limit=%d", category, page+1, limit)
		if status != "" {
			nextURL += fmt.Sprintf("&status=%s", status)
		}
	}
	if page > 1 {
		prevURL = fmt.Sprintf("/api/v1/categories/%s/markets?page=%d&limit=%d", category, page-1, limit)
		if status != "" {
			prevURL += fmt.Sprintf("&status=%s", status)
		}
	}

	pagination := &dto.PaginationDTO{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		NextURL:    nextURL,
		PrevURL:    prevURL,
	}

	return &dto.MarketListDTO{
		Markets:    marketDTOs,
		Pagination: pagination,
	}, nil
}

// marketToDTO converts a market entity to DTO.
func (uc *ListMarkets) marketToDTO(market *entity.Market) *dto.MarketDTO {
	return &dto.MarketDTO{
		Ticker:      market.Ticker.String(),
		Title:       market.Title,
		Category:    market.Category,
		CloseDate:   market.CloseTime,
		YesPrice:    int(market.LastPrice.Value()), // Using LastPrice as yes price for now
		NoPrice:     100 - int(market.LastPrice.Value()),
		Status:      string(market.Status),
		Volume24h:   market.Volume24h,
		LastUpdated: market.LastUpdated,
	}
}
