package repository

import (
	"context"
	"errors"

	"upwork-test/internal/domain/market/entity"
)

var (
	// ErrNotFound is returned when a market is not found
	ErrNotFound = errors.New("market not found")
)

// MarketRepository defines the interface for market data access.
type MarketRepository interface {
	// ListByCategory retrieves markets for a category with pagination
	ListByCategory(ctx context.Context, category string, page int, limit int, status string) ([]*entity.Market, int, error)

	// GetByTicker retrieves a single market by ticker
	GetByTicker(ctx context.Context, ticker string) (*entity.Market, error)

	// GetMultiple retrieves multiple markets by tickers
	GetMultiple(ctx context.Context, tickers []string) ([]*entity.Market, error)

	// GetOrderBook retrieves the order book for a market (cached for 30s)
	GetOrderBook(ctx context.Context, ticker string) (*entity.OrderBook, error)

	// GetRecentTrades retrieves recent trades for a market (cached for 1min)
	GetRecentTrades(ctx context.Context, ticker string, limit int) ([]*entity.Trade, error)
}
