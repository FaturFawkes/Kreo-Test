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
	Ticker      valueobject.Ticker `json:"ticker"`
	Title       string             `json:"title"`
	Category    string             `json:"category"`
	OpenTime    time.Time          `json:"open_time"`
	CloseTime   time.Time          `json:"close_time"`
	Status      MarketStatus       `json:"status"`
	YesAsk      valueobject.Price  `json:"yes_ask"`
	YesBid      valueobject.Price  `json:"yes_bid"`
	NoAsk       valueobject.Price  `json:"no_ask"`
	NoBid       valueobject.Price  `json:"no_bid"`
	LastPrice   valueobject.Price  `json:"last_price"`
	Volume      int64              `json:"volume"`
	Volume24h   int64              `json:"volume_24h"`
	Liquidity   int64              `json:"liquidity"`
	LastUpdated time.Time          `json:"last_updated"`
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
