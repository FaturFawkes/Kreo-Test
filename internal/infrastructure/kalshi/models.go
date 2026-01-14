package kalshi

import "time"

// MarketListResponse represents the response from GET /markets
type MarketListResponse struct {
	Markets []MarketResponse `json:"markets"`
	Cursor  string           `json:"cursor,omitempty"`
}

// MarketResponse represents a market in the API response
type MarketResponse struct {
	Ticker            string    `json:"ticker"`
	EventTicker       string    `json:"event_ticker"`
	Title             string    `json:"title"`
	Subtitle          string    `json:"subtitle"`
	OpenTime          time.Time `json:"open_time"`
	CloseTime         time.Time `json:"close_time"`
	Status            string    `json:"status"`
	Volume            int64     `json:"volume"`
	Volume24h         int64     `json:"volume_24h"`
	Liquidity         int64     `json:"liquidity"`
	YesAsk            int64     `json:"yes_ask"`
	YesBid            int64     `json:"yes_bid"`
	NoAsk             int64     `json:"no_ask"`
	NoBid             int64     `json:"no_bid"`
	LastPrice         int64     `json:"last_price"`
	PreviousYesAsk    int64     `json:"previous_yes_ask"`
	PreviousYesBid    int64     `json:"previous_yes_bid"`
	PreviousPrice     int64     `json:"previous_price"`
	Result            string    `json:"result,omitempty"`
	CanCloseEarly     bool      `json:"can_close_early"`
	ExpirationValue   string    `json:"expiration_value,omitempty"`
	LatestExpiration  time.Time `json:"latest_expiration_time,omitempty"`
	FloorStrike       int64     `json:"floor_strike,omitempty"`
	CapStrike         int64     `json:"cap_strike,omitempty"`
	StrikeType        string    `json:"strike_type,omitempty"`
	SettlementValue   int64     `json:"settlement_value,omitempty"`
	FunctionalStrike  string    `json:"functional_strike,omitempty"`
	RangedGroupTicker string    `json:"ranged_group_ticker,omitempty"`
	Category          string    `json:"-"` // Derived field, not from API
}

// OrderBookResponse represents the order book for a market
type OrderBookResponse struct {
	Ticker     string           `json:"ticker"`
	YesOrders  []OrderBookLevel `json:"yes_orders"`
	NoOrders   []OrderBookLevel `json:"no_orders"`
	LastUpdate time.Time        `json:"last_update"`
}

// OrderBookLevel represents a price level in the order book
type OrderBookLevel struct {
	Price    int64 `json:"price"`
	Quantity int64 `json:"quantity"`
}

// TradesResponse represents the response from GET /markets/{ticker}/trades
type TradesResponse struct {
	Trades []TradeResponse `json:"trades"`
	Cursor string          `json:"cursor,omitempty"`
}

// TradeResponse represents a trade in the API response
type TradeResponse struct {
	TradeID   string    `json:"trade_id"`
	Ticker    string    `json:"ticker"`
	Price     int64     `json:"price"`
	Quantity  int64     `json:"quantity"`
	Side      string    `json:"side"`   // "yes" or "no"
	Action    string    `json:"action"` // "buy" or "sell"
	CreatedAt time.Time `json:"created_at"`
	Taker     string    `json:"taker_side"` // "yes" or "no"
}

// ErrorResponse represents an error response from the Kalshi API
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// SeriesListResponse represents the response from GET /series
type SeriesListResponse struct {
	Series []SeriesResponse `json:"series"`
}

// SeriesResponse represents a series from the Kalshi API
type SeriesResponse struct {
	Ticker    string `json:"ticker"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	Frequency string `json:"frequency"`
}

// EventListResponse represents the response from GET /events
type EventListResponse struct {
	Events []EventResponse `json:"events"`
}

// EventResponse represents an event from the Kalshi API
type EventResponse struct {
	EventTicker  string `json:"event_ticker"`
	SeriesTicker string `json:"series_ticker"`
	Title        string `json:"title"`
	Category     string `json:"category"`
}
