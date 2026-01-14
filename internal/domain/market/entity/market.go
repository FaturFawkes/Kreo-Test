package entity

import (
	"time"
	"upwork-test/internal/domain/market/valueobject"
)

// MarketStatus represents the status of a market
type MarketStatus string

const (
	MarketStatusOpen    MarketStatus = "open"
	MarketStatusClosed  MarketStatus = "closed"
	MarketStatusSettled MarketStatus = "settled"
)

// Market represents a prediction market
type Market struct {
	Ticker      valueobject.Ticker
	Title       string
	Category    string
	OpenTime    time.Time
	CloseTime   time.Time
	Status      MarketStatus
	YesAsk      valueobject.Price
	YesBid      valueobject.Price
	NoAsk       valueobject.Price
	NoBid       valueobject.Price
	LastPrice   valueobject.Price
	Volume      int64
	Volume24h   int64
	Liquidity   int64
	LastUpdated time.Time
}

// NewMarket creates a new Market entity
func NewMarket(
	ticker valueobject.Ticker,
	title string,
	category string,
	openTime time.Time,
	closeTime time.Time,
	status MarketStatus,
) *Market {
	return &Market{
		Ticker:      ticker,
		Title:       title,
		Category:    category,
		OpenTime:    openTime,
		CloseTime:   closeTime,
		Status:      status,
		LastUpdated: time.Now(),
	}
}

// IsOpen checks if the market is currently open
func (m *Market) IsOpen() bool {
	return m.Status == MarketStatusOpen
}

// IsClosed checks if the market is closed
func (m *Market) IsClosed() bool {
	return m.Status == MarketStatusClosed || m.Status == MarketStatusSettled
}
