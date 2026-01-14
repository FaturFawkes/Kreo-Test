package entity

import (
	"time"

	"upwork-test/internal/domain/market/valueobject"
)

// OrderLevel represents a price level in the order book
type OrderLevel struct {
	Price    valueobject.Price
	Quantity int
}

// NewOrderLevel creates a new OrderLevel
func NewOrderLevel(price valueobject.Price, quantity int) *OrderLevel {
	return &OrderLevel{
		Price:    price,
		Quantity: quantity,
	}
}

// OrderBook represents a snapshot of bids and asks for a market
type OrderBook struct {
	Ticker    valueobject.Ticker
	Timestamp time.Time
	Bids      []OrderLevel
	Asks      []OrderLevel
}

// NewOrderBook creates a new OrderBook entity
func NewOrderBook(
	ticker valueobject.Ticker,
	timestamp time.Time,
	bids []OrderLevel,
	asks []OrderLevel,
) *OrderBook {
	return &OrderBook{
		Ticker:    ticker,
		Timestamp: timestamp,
		Bids:      bids,
		Asks:      asks,
	}
}

// Spread returns the difference between best bid and best ask in cents
func (ob *OrderBook) Spread() int64 {
	if len(ob.Bids) == 0 || len(ob.Asks) == 0 {
		return 0
	}

	bestBid := ob.Bids[0].Price.Value()
	bestAsk := ob.Asks[0].Price.Value()

	if bestAsk > bestBid {
		return bestAsk - bestBid
	}

	return 0
}

// BidDepth returns total quantity across all bid levels
func (ob *OrderBook) BidDepth() int {
	depth := 0
	for _, level := range ob.Bids {
		depth += level.Quantity
	}
	return depth
}

// AskDepth returns total quantity across all ask levels
func (ob *OrderBook) AskDepth() int {
	depth := 0
	for _, level := range ob.Asks {
		depth += level.Quantity
	}
	return depth
}

// TotalDepth returns total quantity across all bid and ask levels
func (ob *OrderBook) TotalDepth() int {
	return ob.BidDepth() + ob.AskDepth()
}

// BestBid returns the highest bid price, or nil if no bids
func (ob *OrderBook) BestBid() *valueobject.Price {
	if len(ob.Bids) == 0 {
		return nil
	}
	price := ob.Bids[0].Price
	return &price
}

// BestAsk returns the lowest ask price, or nil if no asks
func (ob *OrderBook) BestAsk() *valueobject.Price {
	if len(ob.Asks) == 0 {
		return nil
	}
	price := ob.Asks[0].Price
	return &price
}

// IsEmpty returns true if there are no bids or asks
func (ob *OrderBook) IsEmpty() bool {
	return len(ob.Bids) == 0 && len(ob.Asks) == 0
}
