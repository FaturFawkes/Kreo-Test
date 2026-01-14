package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidCategoryName is returned when category name is invalid
	ErrInvalidCategoryName = errors.New("invalid category name")

	// validCategories defines the allowed category names from Kalshi
	validCategories = map[string]bool{
		"ECONOMICS":     true,
		"POLITICS":      true,
		"SPORTS":        true,
		"CLIMATE":       true,
		"SCIENCE":       true,
		"TECHNOLOGY":    true,
		"ENTERTAINMENT": true,
		"FINANCE":       true,
		"HEALTH":        true,
		"CRYPTO":        true,
	}
)

// CategoryName represents a market category name
type CategoryName struct {
	value string
}

// NewCategoryName creates a new CategoryName value object
func NewCategoryName(value string) (CategoryName, error) {
	// Normalize: trim whitespace and convert to uppercase
	normalized := strings.ToUpper(strings.TrimSpace(value))

	// Validate
	if normalized == "" {
		return CategoryName{}, fmt.Errorf("%w: category name cannot be empty", ErrInvalidCategoryName)
	}

	if !validCategories[normalized] {
		return CategoryName{}, fmt.Errorf("%w: category '%s' is not a valid Kalshi category", ErrInvalidCategoryName, normalized)
	}

	return CategoryName{value: normalized}, nil
}

// String returns the string representation of the category name
func (c CategoryName) String() string {
	return c.value
}

// Equals checks if two category names are equal
func (c CategoryName) Equals(other CategoryName) bool {
	return c.value == other.value
}

// IsEmpty checks if the category name is empty
func (c CategoryName) IsEmpty() bool {
	return c.value == ""
}

// IsValid checks if a category name string is valid (without creating an object)
func IsValid(value string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	return validCategories[normalized]
}

// AllCategories returns a list of all valid category names
func AllCategories() []string {
	categories := make([]string, 0, len(validCategories))
	for category := range validCategories {
		categories = append(categories, category)
	}
	return categories
}

// MarshalJSON implements json.Marshaler
func (c CategoryName) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.value + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (c *CategoryName) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return ErrInvalidCategoryName
	}
	// Remove quotes
	value := string(data[1 : len(data)-1])
	c.value = value
	return nil
}
