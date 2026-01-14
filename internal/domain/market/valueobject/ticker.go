package valueobject

import (
	"errors"
	"regexp"
)

var (
	// ErrInvalidTicker is returned when ticker format is invalid
	ErrInvalidTicker = errors.New("invalid ticker format")

	// tickerPattern defines the valid ticker format
	// Kalshi tickers are uppercase alphanumeric with hyphens
	tickerPattern = regexp.MustCompile(`^[A-Z0-9\-]+$`)
)

// Ticker represents a market ticker symbol
type Ticker struct {
	value string
}

// NewTicker creates a new Ticker value object
func NewTicker(value string) (Ticker, error) {
	return Ticker{value: value}, nil
}

// String returns the string representation of the ticker
func (t Ticker) String() string {
	return t.value
}

// Equals checks if two tickers are equal
func (t Ticker) Equals(other Ticker) bool {
	return t.value == other.value
}

// IsEmpty checks if the ticker is empty
func (t Ticker) IsEmpty() bool {
	return t.value == ""
}
