package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"upwork-test/internal/domain/market/entity"
	"upwork-test/internal/domain/market/valueobject"
	"upwork-test/internal/infrastructure/kalshi"

	"github.com/redis/go-redis/v9"
)


const (
	marketListCacheTTL = 5 * time.Minute
)

// MarketRepository implements the market repository with Redis caching.
type MarketRepository struct {
	redisClient  *redis.Client
	kalshiClient *kalshi.Client
	keyBuilder   *KeyBuilder
	mapper       *kalshi.Mapper
}

// NewMarketRepository creates a new market repository.
func NewMarketRepository(redisClient *redis.Client, kalshiClient *kalshi.Client) *MarketRepository {
	return &MarketRepository{
		redisClient:  redisClient,
		kalshiClient: kalshiClient,
		keyBuilder:   NewKeyBuilder("kalshi"),
		mapper:       kalshi.NewMapper(),
	}
}

// ListByCategory retrieves markets for a category with pagination.
func (r *MarketRepository) ListByCategory(ctx context.Context, category string, page int, limit int, status string) ([]*entity.Market, int, error) {
	cacheKey := r.keyBuilder.MarketList(category)

	cachedData, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var markets []*entity.Market
		if err := json.Unmarshal([]byte(cachedData), &markets); err == nil {
			filtered := r.filterByStatus(markets, status)
			paginated, total := r.paginate(filtered, page, limit)
			return paginated, total, nil
		}
	}

	kalshiResponse, err := r.kalshiClient.GetMarkets(ctx, category, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch markets from Kalshi: %w", err)
	}

	markets, err := r.mapper.ToMarketEntities(kalshiResponse.Markets)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to map markets: %w", err)
	}

	if data, err := json.Marshal(markets); err == nil {
		r.redisClient.Set(ctx, cacheKey, data, marketListCacheTTL)
	}

	filtered := r.filterByStatus(markets, status)
	paginated, total := r.paginate(filtered, page, limit)

	return paginated, total, nil
}

// GetByTicker retrieves a single market by ticker.
func (r *MarketRepository) GetByTicker(ctx context.Context, tickerStr string) (*entity.Market, error) {
	cacheKey := r.keyBuilder.MarketMetadata(tickerStr)

	cachedData, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var market entity.Market
		if err := json.Unmarshal([]byte(cachedData), &market); err == nil {
			return &market, nil
		}
	}

	ticker, err := valueobject.NewTicker(tickerStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ticker: %w", err)
	}

	kalshiMarket, err := r.kalshiClient.GetMarket(ctx, ticker.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market from Kalshi: %w", err)
	}

	market, err := r.mapper.ToMarketEntity(kalshiMarket)
	if err != nil {
		return nil, fmt.Errorf("failed to map market: %w", err)
	}

	if data, err := json.Marshal(market); err == nil {
		r.redisClient.Set(ctx, cacheKey, data, marketListCacheTTL)
	}

	return market, nil
}

// GetMultiple retrieves multiple markets by tickers.
func (r *MarketRepository) GetMultiple(ctx context.Context, tickers []string) ([]*entity.Market, error) {
	markets := make([]*entity.Market, 0, len(tickers))

	for _, ticker := range tickers {
		market, err := r.GetByTicker(ctx, ticker)
		if err != nil {
			// Skip markets that fail to fetch
			continue
		}
		markets = append(markets, market)
	}

	return markets, nil
}

// filterByStatus filters markets by status.
func (r *MarketRepository) filterByStatus(markets []*entity.Market, status string) []*entity.Market {
	if status == "" {
		return markets
	}

	filtered := make([]*entity.Market, 0, len(markets))
	normalizedStatus := strings.ToLower(status)

	for _, market := range markets {
		if strings.ToLower(string(market.Status)) == normalizedStatus {
			filtered = append(filtered, market)
		}
	}

	return filtered
}

// paginate paginates the market list.
func (r *MarketRepository) paginate(markets []*entity.Market, page int, limit int) ([]*entity.Market, int) {
	total := len(markets)

	start := (page - 1) * limit
	if start >= total {
		return []*entity.Market{}, total
	}

	end := start + limit
	if end > total {
		end = total
	}

	return markets[start:end], total
}

// GetOrderBook retrieves the order book for a market (cached for 30s)
func (r *MarketRepository) GetOrderBook(ctx context.Context, ticker string) (*entity.OrderBook, error) {
	cacheKey := r.keyBuilder.MarketOrderBook(ticker)

	cachedData, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var orderBook entity.OrderBook
		if err := json.Unmarshal([]byte(cachedData), &orderBook); err == nil {
			return &orderBook, nil
		}
	}

	kalshiResponse, err := r.kalshiClient.GetOrderBook(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order book from Kalshi: %w", err)
	}

	orderBook, err := r.mapper.ToOrderBookEntity(kalshiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to convert order book: %w", err)
	}

	if data, err := json.Marshal(orderBook); err == nil {
		_ = r.redisClient.Set(ctx, cacheKey, data, 30*time.Second).Err()
	}

	return orderBook, nil
}

// GetRecentTrades retrieves recent trades for a market (cached for 1min)
func (r *MarketRepository) GetRecentTrades(ctx context.Context, ticker string, limit int) ([]*entity.Trade, error) {
	cacheKey := r.keyBuilder.MarketTrades(ticker)

	cachedData, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var trades []*entity.Trade
		if err := json.Unmarshal([]byte(cachedData), &trades); err == nil {
			return trades, nil
		}
	}

	kalshiResponse, err := r.kalshiClient.GetTrades(ctx, ticker, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trades from Kalshi: %w", err)
	}

	trades, err := r.mapper.ToTradeEntities(kalshiResponse.Trades)
	if err != nil {
		return nil, fmt.Errorf("failed to convert trades: %w", err)
	}

	if data, err := json.Marshal(trades); err == nil {
		_ = r.redisClient.Set(ctx, cacheKey, data, 1*time.Minute).Err()
	}

	return trades, nil
}
