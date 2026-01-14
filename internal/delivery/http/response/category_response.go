package response

import (
	"time"
	"upwork-test/internal/application/dto"
)

// CategoryOverviewResponse represents category overview metrics in the API response.
type CategoryOverviewResponse struct {
	CategoryName     string    `json:"category_name"`
	TotalMarkets     int       `json:"total_markets"`
	TotalVolume24h   int64     `json:"total_volume_24h"`
	AverageLiquidity float64   `json:"average_liquidity"`
	ActiveTraders24h int       `json:"active_traders_24h"`
	ComputedAt       time.Time `json:"computed_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}

// FromCategoryOverviewDTO converts a category overview DTO to API response format.
func FromCategoryOverviewDTO(overviewDTO *dto.CategoryOverviewDTO) *CategoryOverviewResponse {
	return &CategoryOverviewResponse{
		CategoryName:     overviewDTO.CategoryName,
		TotalMarkets:     overviewDTO.TotalMarkets,
		TotalVolume24h:   overviewDTO.TotalVolume24h,
		AverageLiquidity: overviewDTO.AverageLiquidity,
		ActiveTraders24h: overviewDTO.ActiveTraders24h,
		ComputedAt:       overviewDTO.ComputedAt,
		ExpiresAt:        overviewDTO.ExpiresAt,
	}
}
