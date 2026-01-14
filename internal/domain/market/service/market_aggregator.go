package service

import (
	"upwork-test/internal/domain/market/entity"
)

// AggregatedMarket represents a market with all its data combined
type AggregatedMarket struct {
	Market    *entity.Market
	OrderBook *entity.OrderBook
	Trades    []*entity.Trade
	IsPartial bool // True if some components failed to load
}

// MarketAggregator combines data from multiple sources into a complete market view
type MarketAggregator struct{}

// NewMarketAggregator creates a new MarketAggregator service
func NewMarketAggregator() *MarketAggregator {
	return &MarketAggregator{}
}

// Aggregate combines market metadata, order book, and trades into a single view
// It returns a partial result if some components are missing
func (ma *MarketAggregator) Aggregate(
	market *entity.Market,
	orderBook *entity.OrderBook,
	trades []*entity.Trade,
) *AggregatedMarket {
	isPartial := false

	// If market is nil, we cannot aggregate anything
	if market == nil {
		return nil
	}

	// Check if we're missing optional components
	if orderBook == nil || trades == nil || len(trades) == 0 {
		isPartial = true
	}

	return &AggregatedMarket{
		Market:    market,
		OrderBook: orderBook,
		Trades:    trades,
		IsPartial: isPartial,
	}
}

// HasOrderBook returns true if the aggregated market includes an order book
func (am *AggregatedMarket) HasOrderBook() bool {
	return am.OrderBook != nil
}

// HasTrades returns true if the aggregated market includes trades
func (am *AggregatedMarket) HasTrades() bool {
	return am.Trades != nil && len(am.Trades) > 0
}

// IsComplete returns true if all components are present
func (am *AggregatedMarket) IsComplete() bool {
	return !am.IsPartial
}
