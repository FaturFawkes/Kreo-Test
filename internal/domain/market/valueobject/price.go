package valueobject

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidPrice is returned when price is invalid
	ErrInvalidPrice = errors.New("invalid price")
)

// Price represents a market price in cents (0-100)
type Price struct {
	value int64
}

// NewPrice creates a new Price value object
func NewPrice(value int64) (Price, error) {
	// Validate range (Kalshi prices are 0-100 cents)
	if value < 0 {
		return Price{}, fmt.Errorf("%w: price cannot be negative", ErrInvalidPrice)
	}

	if value > 100 {
		return Price{}, fmt.Errorf("%w: price cannot exceed 100 cents", ErrInvalidPrice)
	}

	return Price{value: value}, nil
}

// Value returns the price value in cents
func (p Price) Value() int64 {
	return p.value
}

// Cents returns the price in cents (same as Value)
func (p Price) Cents() int64 {
	return p.value
}

// Dollars returns the price in dollars as a float
func (p Price) Dollars() float64 {
	return float64(p.value) / 100.0
}

// Equals checks if two prices are equal
func (p Price) Equals(other Price) bool {
	return p.value == other.value
}

// MarshalJSON implements json.Marshaler
func (p Price) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", p.value)), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (p *Price) UnmarshalJSON(data []byte) error {
	var value int64
	if _, err := fmt.Sscanf(string(data), "%d", &value); err != nil {
		return err
	}
	p.value = value
	return nil
}

// GreaterThan checks if this price is greater than another
func (p Price) GreaterThan(other Price) bool {
	return p.value > other.value
}

// LessThan checks if this price is less than another
func (p Price) LessThan(other Price) bool {
	return p.value < other.value
}

// IsZero checks if the price is zero
func (p Price) IsZero() bool {
	return p.value == 0
}
