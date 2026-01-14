package entity

import (
	"errors"
	"time"

	"upwork-test/internal/domain/category/valueobject"
)

var (
	// ErrInvalidTotalMarkets is returned when total markets is negative
	ErrInvalidTotalMarkets = errors.New("total markets cannot be negative")
	// ErrInvalidTotalVolume is returned when total volume is negative
	ErrInvalidTotalVolume = errors.New("total volume cannot be negative")
	// ErrInvalidAverageLiquidity is returned when average liquidity is negative
	ErrInvalidAverageLiquidity = errors.New("average liquidity cannot be negative")
	// ErrInvalidActiveTraders is returned when active traders is negative
	ErrInvalidActiveTraders = errors.New("active traders cannot be negative")
	// ErrInvalidTimeRange is returned when computed time is after expiry time
	ErrInvalidTimeRange = errors.New("computed time must be before expiry time")
)

// CategoryOverview represents pre-computed aggregate metrics for a category
type CategoryOverview struct {
	CategoryName     valueobject.CategoryName `json:"category_name"`
	TotalMarkets     int                      `json:"total_markets"`
	TotalVolume24h   int64                    `json:"total_volume_24h"`
	AverageLiquidity float64                  `json:"average_liquidity"`
	ActiveTraders24h int                      `json:"active_traders_24h"`
	ComputedAt       time.Time                `json:"computed_at"`
	ExpiresAt        time.Time                `json:"expires_at"`
}

// NewCategoryOverview creates a new CategoryOverview entity
func NewCategoryOverview(
	categoryName valueobject.CategoryName,
	totalMarkets int,
	totalVolume24h int64,
	averageLiquidity float64,
	activeTraders24h int,
	ttl time.Duration,
) (*CategoryOverview, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)

	// Validate invariants
	if totalMarkets < 0 {
		return nil, ErrInvalidTotalMarkets
	}
	if totalVolume24h < 0 {
		return nil, ErrInvalidTotalVolume
	}
	if averageLiquidity < 0 {
		return nil, ErrInvalidAverageLiquidity
	}
	if activeTraders24h < 0 {
		return nil, ErrInvalidActiveTraders
	}

	return &CategoryOverview{
		CategoryName:     categoryName,
		TotalMarkets:     totalMarkets,
		TotalVolume24h:   totalVolume24h,
		AverageLiquidity: averageLiquidity,
		ActiveTraders24h: activeTraders24h,
		ComputedAt:       now,
		ExpiresAt:        expiresAt,
	}, nil
}

// IsExpired checks if the overview data has expired
func (co *CategoryOverview) IsExpired() bool {
	return time.Now().After(co.ExpiresAt)
}

// IsFresh checks if the overview data is still fresh (not expired)
func (co *CategoryOverview) IsFresh() bool {
	return !co.IsExpired()
}

// TimeUntilExpiry returns the duration until this overview expires
func (co *CategoryOverview) TimeUntilExpiry() time.Duration {
	if co.IsExpired() {
		return 0
	}
	return time.Until(co.ExpiresAt)
}

// Age returns how long ago this overview was computed
func (co *CategoryOverview) Age() time.Duration {
	return time.Since(co.ComputedAt)
}
