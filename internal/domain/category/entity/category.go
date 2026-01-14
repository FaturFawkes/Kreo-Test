package entity

import (
	"errors"
	"time"

	"upwork-test/internal/domain/category/valueobject"
)

var (
	// ErrInvalidMarketCount is returned when market count is negative
	ErrInvalidMarketCount = errors.New("market count cannot be negative")
)

// Category represents a logical grouping of markets
type Category struct {
	Name        valueobject.CategoryName `json:"name"`
	DisplayName string                   `json:"display_name"`
	Description string                   `json:"description"`
	MarketCount int                      `json:"market_count"`
	LastUpdated time.Time                `json:"last_updated"`
}

// NewCategory creates a new Category entity
func NewCategory(
	name valueobject.CategoryName,
	displayName string,
	description string,
) *Category {
	return &Category{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		MarketCount: 0,
		LastUpdated: time.Now(),
	}
}

// UpdateMarketCount updates the market count for this category
func (c *Category) UpdateMarketCount(count int) error {
	if count < 0 {
		return ErrInvalidMarketCount
	}
	c.MarketCount = count
	c.LastUpdated = time.Now()
	return nil
}

// IsEmpty checks if the category has no markets
func (c *Category) IsEmpty() bool {
	return c.MarketCount == 0
}
