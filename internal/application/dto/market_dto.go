package dto

import (
	"time"
)

// MarketDTO represents a market data transfer object for the application layer.
type MarketDTO struct {
	Ticker          string    `json:"ticker"`
	Title           string    `json:"title"`
	Category        string    `json:"category"`
	CloseDate       time.Time `json:"close_date"`
	SettlementRules string    `json:"settlement_rules,omitempty"`
	YesPrice        int       `json:"yes_price"`
	NoPrice         int       `json:"no_price"`
	Status          string    `json:"status"`
	Volume24h       int64     `json:"volume_24h"`
	LastUpdated     time.Time `json:"last_updated"`
}

// PaginationDTO represents pagination metadata.
type PaginationDTO struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
	NextURL    string `json:"next_url,omitempty"`
	PrevURL    string `json:"prev_url,omitempty"`
}

// MarketListDTO represents a paginated list of markets.
type MarketListDTO struct {
	Markets    []*MarketDTO   `json:"data"`
	Pagination *PaginationDTO `json:"pagination"`
}

// MarketDetailDTO represents comprehensive market information with aggregated data
type MarketDetailDTO struct {
	Ticker       string        `json:"ticker"`
	Title        string        `json:"title"`
	Category     string        `json:"category"`
	OpenTime     time.Time     `json:"open_time"`
	CloseTime    time.Time     `json:"close_time"`
	Status       string        `json:"status"`
	YesAsk       int64         `json:"yes_ask"`
	YesBid       int64         `json:"yes_bid"`
	NoAsk        int64         `json:"no_ask"`
	NoBid        int64         `json:"no_bid"`
	LastPrice    int64         `json:"last_price"`
	Volume       int64         `json:"volume"`
	Volume24h    int64         `json:"volume_24h"`
	Liquidity    int64         `json:"liquidity"`
	OrderBook    *OrderBookDTO `json:"order_book,omitempty"`
	RecentTrades []TradeDTO    `json:"recent_trades,omitempty"`
	IsPartial    bool          `json:"is_partial"`
	Errors       []string      `json:"errors,omitempty"`
}

// OrderBookDTO represents an order book snapshot
type OrderBookDTO struct {
	Timestamp time.Time       `json:"timestamp"`
	Bids      []OrderLevelDTO `json:"bids"`
	Asks      []OrderLevelDTO `json:"asks"`
	Spread    int64           `json:"spread"`
}

// OrderLevelDTO represents a price level in the order book
type OrderLevelDTO struct {
	Price    int64 `json:"price"`
	Quantity int   `json:"quantity"`
}

// TradeDTO represents a historical trade
type TradeDTO struct {
	TradeID   string    `json:"trade_id"`
	Price     int64     `json:"price"`
	Quantity  int       `json:"quantity"`
	Side      string    `json:"side"`
	Timestamp time.Time `json:"timestamp"`
}
