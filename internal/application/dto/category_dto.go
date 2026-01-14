package dto

import (
	"time"
)

// CategoryDTO represents a category data transfer object for the application layer.
type CategoryDTO struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	MarketCount int       `json:"market_count"`
	LastUpdated time.Time `json:"last_updated"`
}

// CategoryOverviewDTO represents category overview metrics.
type CategoryOverviewDTO struct {
	CategoryName     string    `json:"category_name"`
	TotalMarkets     int       `json:"total_markets"`
	TotalVolume24h   int64     `json:"total_volume_24h"`
	AverageLiquidity float64   `json:"average_liquidity"`
	ActiveTraders24h int       `json:"active_traders_24h"`
	ComputedAt       time.Time `json:"computed_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}
