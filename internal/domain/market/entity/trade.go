package entity

import (
	"time"

	"upwork-test/internal/domain/market/valueobject"
)

// TradeSide represents the side of a trade (buy or sell)
type TradeSide string

const (
	TradeSideBuy  TradeSide = "buy"
	TradeSideSell TradeSide = "sell"
)

// Trade represents a historical trade record for a market
type Trade struct {
	Ticker    valueobject.Ticker
	TradeID   string
	Price     valueobject.Price
	Quantity  int
	Side      TradeSide
	Timestamp time.Time
}

// NewTrade creates a new Trade entity
func NewTrade(
	ticker valueobject.Ticker,
	tradeID string,
	price valueobject.Price,
	quantity int,
	side TradeSide,
	timestamp time.Time,
) *Trade {
	return &Trade{
		Ticker:    ticker,
		TradeID:   tradeID,
		Price:     price,
		Quantity:  quantity,
		Side:      side,
		Timestamp: timestamp,
	}
}

// IsBuy returns true if this trade is a buy
func (t *Trade) IsBuy() bool {
	return t.Side == TradeSideBuy
}

// IsSell returns true if this trade is a sell
func (t *Trade) IsSell() bool {
	return t.Side == TradeSideSell
}

// Value returns the total value of the trade in cents (price * quantity)
func (t *Trade) Value() int64 {
	return t.Price.Value() * int64(t.Quantity)
}
