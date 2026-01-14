package response

import (
	"time"
	"upwork-test/internal/application/dto"
)

// MarketResponse represents a market in the API response.
type MarketResponse struct {
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

// PaginationResponse represents pagination metadata in the API response.
type PaginationResponse struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
	NextURL    string `json:"next_url,omitempty"`
	PrevURL    string `json:"prev_url,omitempty"`
}

// MarketListResponse represents the response for market listing.
type MarketListResponse struct {
	Data       []*MarketResponse   `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

// FromMarketDTO converts a market DTO to API response format.
func FromMarketDTO(marketDTO *dto.MarketDTO) *MarketResponse {
	return &MarketResponse{
		Ticker:          marketDTO.Ticker,
		Title:           marketDTO.Title,
		Category:        marketDTO.Category,
		CloseDate:       marketDTO.CloseDate,
		SettlementRules: marketDTO.SettlementRules,
		YesPrice:        marketDTO.YesPrice,
		NoPrice:         marketDTO.NoPrice,
		Status:          marketDTO.Status,
		Volume24h:       marketDTO.Volume24h,
		LastUpdated:     marketDTO.LastUpdated,
	}
}

// FromMarketListDTO converts a market list DTO to API response format.
func FromMarketListDTO(listDTO *dto.MarketListDTO) *MarketListResponse {
	markets := make([]*MarketResponse, len(listDTO.Markets))
	for i, marketDTO := range listDTO.Markets {
		markets[i] = FromMarketDTO(marketDTO)
	}

	pagination := &PaginationResponse{
		Page:       listDTO.Pagination.Page,
		Limit:      listDTO.Pagination.Limit,
		Total:      listDTO.Pagination.Total,
		TotalPages: listDTO.Pagination.TotalPages,
		NextURL:    listDTO.Pagination.NextURL,
		PrevURL:    listDTO.Pagination.PrevURL,
	}

	return &MarketListResponse{
		Data:       markets,
		Pagination: pagination,
	}
}

// MarketDetailResponse represents comprehensive market information with aggregated data
type MarketDetailResponse struct {
	Ticker       string             `json:"ticker"`
	Title        string             `json:"title"`
	Category     string             `json:"category"`
	OpenTime     time.Time          `json:"open_time"`
	CloseTime    time.Time          `json:"close_time"`
	Status       string             `json:"status"`
	YesAsk       int64              `json:"yes_ask"`
	YesBid       int64              `json:"yes_bid"`
	NoAsk        int64              `json:"no_ask"`
	NoBid        int64              `json:"no_bid"`
	LastPrice    int64              `json:"last_price"`
	Volume       int64              `json:"volume"`
	Volume24h    int64              `json:"volume_24h"`
	Liquidity    int64              `json:"liquidity"`
	OrderBook    *OrderBookResponse `json:"order_book,omitempty"`
	RecentTrades []TradeResponse    `json:"recent_trades,omitempty"`
	IsPartial    bool               `json:"is_partial"`
	Errors       []string           `json:"errors,omitempty"`
}

// OrderBookResponse represents an order book snapshot
type OrderBookResponse struct {
	Timestamp time.Time            `json:"timestamp"`
	Bids      []OrderLevelResponse `json:"bids"`
	Asks      []OrderLevelResponse `json:"asks"`
	Spread    int64                `json:"spread"`
}

// OrderLevelResponse represents a price level in the order book
type OrderLevelResponse struct {
	Price    int64 `json:"price"`
	Quantity int   `json:"quantity"`
}

// TradeResponse represents a historical trade
type TradeResponse struct {
	TradeID   string    `json:"trade_id"`
	Price     int64     `json:"price"`
	Quantity  int       `json:"quantity"`
	Side      string    `json:"side"`
	Timestamp time.Time `json:"timestamp"`
}

// FromMarketDetailDTO converts a market detail DTO to API response format
func FromMarketDetailDTO(detailDTO *dto.MarketDetailDTO) *MarketDetailResponse {
	response := &MarketDetailResponse{
		Ticker:    detailDTO.Ticker,
		Title:     detailDTO.Title,
		Category:  detailDTO.Category,
		OpenTime:  detailDTO.OpenTime,
		CloseTime: detailDTO.CloseTime,
		Status:    detailDTO.Status,
		YesAsk:    detailDTO.YesAsk,
		YesBid:    detailDTO.YesBid,
		NoAsk:     detailDTO.NoAsk,
		NoBid:     detailDTO.NoBid,
		LastPrice: detailDTO.LastPrice,
		Volume:    detailDTO.Volume,
		Volume24h: detailDTO.Volume24h,
		Liquidity: detailDTO.Liquidity,
		IsPartial: detailDTO.IsPartial,
		Errors:    detailDTO.Errors,
	}

	// Convert order book if present
	if detailDTO.OrderBook != nil {
		response.OrderBook = &OrderBookResponse{
			Timestamp: detailDTO.OrderBook.Timestamp,
			Bids:      convertOrderLevels(detailDTO.OrderBook.Bids),
			Asks:      convertOrderLevels(detailDTO.OrderBook.Asks),
			Spread:    detailDTO.OrderBook.Spread,
		}
	}

	// Convert trades if present
	if len(detailDTO.RecentTrades) > 0 {
		response.RecentTrades = convertTrades(detailDTO.RecentTrades)
	}

	return response
}

// convertOrderLevels converts DTO order levels to response format
func convertOrderLevels(levels []dto.OrderLevelDTO) []OrderLevelResponse {
	result := make([]OrderLevelResponse, len(levels))
	for i, level := range levels {
		result[i] = OrderLevelResponse{
			Price:    level.Price,
			Quantity: level.Quantity,
		}
	}
	return result
}

// convertTrades converts DTO trades to response format
func convertTrades(trades []dto.TradeDTO) []TradeResponse {
	result := make([]TradeResponse, len(trades))
	for i, trade := range trades {
		result[i] = TradeResponse{
			TradeID:   trade.TradeID,
			Price:     trade.Price,
			Quantity:  trade.Quantity,
			Side:      trade.Side,
			Timestamp: trade.Timestamp,
		}
	}
	return result
}
