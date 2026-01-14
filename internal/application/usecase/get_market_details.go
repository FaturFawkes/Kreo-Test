package usecase

import (
	"context"
	"errors"
	"sync"
	"upwork-test/internal/application/dto"
	"upwork-test/internal/domain/market/entity"
	"upwork-test/internal/domain/market/repository"
	"upwork-test/internal/domain/market/service"
	"upwork-test/internal/domain/market/valueobject"
)

var (
	// ErrMarketNotFound is returned when the market does not exist
	ErrMarketNotFound = errors.New("market not found")
	// ErrInvalidTicker is returned when the ticker is invalid
	ErrInvalidTicker = errors.New("invalid ticker")
)

// MarketRepositoryExtended extends MarketRepository with aggregation methods
type MarketRepositoryExtended interface {
	repository.MarketRepository

	// GetOrderBook retrieves the order book for a market
	GetOrderBook(ctx context.Context, ticker string) (*entity.OrderBook, error)

	// GetRecentTrades retrieves recent trades for a market
	GetRecentTrades(ctx context.Context, ticker string, limit int) ([]*entity.Trade, error)
}

// GetMarketDetails retrieves comprehensive market information with concurrent aggregation
type GetMarketDetails struct {
	repo       MarketRepositoryExtended
	aggregator *service.MarketAggregator
}

// NewGetMarketDetails creates a new GetMarketDetails use case
func NewGetMarketDetails(repo MarketRepositoryExtended) *GetMarketDetails {
	return &GetMarketDetails{
		repo:       repo,
		aggregator: service.NewMarketAggregator(),
	}
}

// Execute retrieves market details with concurrent fan-out for order book and trades
func (uc *GetMarketDetails) Execute(ctx context.Context, tickerStr string) (*dto.MarketDetailDTO, error) {
	ticker, err := valueobject.NewTicker(tickerStr)
	if err != nil {
		return nil, ErrInvalidTicker
	}

	var wg sync.WaitGroup
	var market *entity.Market
	var orderBook *entity.OrderBook
	var trades []*entity.Trade
	var marketErr, orderBookErr, tradesErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		market, marketErr = uc.repo.GetByTicker(ctx, ticker.String())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		orderBook, orderBookErr = uc.repo.GetOrderBook(ctx, ticker.String())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		trades, tradesErr = uc.repo.GetRecentTrades(ctx, ticker.String(), 100)
	}()

	wg.Wait()

	if marketErr != nil {
		if errors.Is(marketErr, repository.ErrNotFound) {
			return nil, ErrMarketNotFound
		}
		return nil, marketErr
	}

	aggregated := uc.aggregator.Aggregate(market, orderBook, trades)
	if aggregated == nil {
		return nil, ErrMarketNotFound
	}

	return uc.toDTO(aggregated, orderBookErr, tradesErr), nil
}

// toDTO converts aggregated market to DTO
func (uc *GetMarketDetails) toDTO(
	aggregated *service.AggregatedMarket,
	orderBookErr, tradesErr error,
) *dto.MarketDetailDTO {
	market := aggregated.Market

	result := &dto.MarketDetailDTO{
		Ticker:    market.Ticker.String(),
		Title:     market.Title,
		Category:  market.Category,
		OpenTime:  market.OpenTime,
		CloseTime: market.CloseTime,
		Status:    string(market.Status),
		YesAsk:    market.YesAsk.Value(),
		YesBid:    market.YesBid.Value(),
		NoAsk:     market.NoAsk.Value(),
		NoBid:     market.NoBid.Value(),
		LastPrice: market.LastPrice.Value(),
		Volume:    market.Volume,
		Volume24h: market.Volume24h,
		Liquidity: market.Liquidity,
		IsPartial: aggregated.IsPartial,
	}

	if aggregated.HasOrderBook() {
		ob := aggregated.OrderBook
		result.OrderBook = &dto.OrderBookDTO{
			Timestamp: ob.Timestamp,
			Bids:      uc.convertOrderLevels(ob.Bids),
			Asks:      uc.convertOrderLevels(ob.Asks),
			Spread:    ob.Spread(),
		}
	} else if orderBookErr != nil {
		result.Errors = append(result.Errors, "order_book: "+orderBookErr.Error())
	}

	if aggregated.HasTrades() {
		result.RecentTrades = uc.convertTrades(aggregated.Trades)
	} else if tradesErr != nil {
		result.Errors = append(result.Errors, "trades: "+tradesErr.Error())
	}

	return result
}

// convertOrderLevels converts domain order levels to DTOs
func (uc *GetMarketDetails) convertOrderLevels(levels []entity.OrderLevel) []dto.OrderLevelDTO {
	result := make([]dto.OrderLevelDTO, len(levels))
	for i, level := range levels {
		result[i] = dto.OrderLevelDTO{
			Price:    level.Price.Value(),
			Quantity: level.Quantity,
		}
	}
	return result
}

// convertTrades converts domain trades to DTOs
func (uc *GetMarketDetails) convertTrades(trades []*entity.Trade) []dto.TradeDTO {
	result := make([]dto.TradeDTO, len(trades))
	for i, trade := range trades {
		result[i] = dto.TradeDTO{
			TradeID:   trade.TradeID,
			Price:     trade.Price.Value(),
			Quantity:  trade.Quantity,
			Side:      string(trade.Side),
			Timestamp: trade.Timestamp,
		}
	}
	return result
}
