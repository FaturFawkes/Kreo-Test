package kalshi

import (
	"fmt"
	"upwork-test/internal/domain/market/entity"
	"upwork-test/internal/domain/market/valueobject"
)

// Mapper converts Kalshi API models to domain entities
type Mapper struct{}

// NewMapper creates a new Mapper
func NewMapper() *Mapper {
	return &Mapper{}
}

// ToMarketEntity converts a MarketResponse to a Market entity
func (m *Mapper) ToMarketEntity(resp *MarketResponse) (*entity.Market, error) {
	ticker, err := valueobject.NewTicker(resp.Ticker)
	if err != nil {
		return nil, fmt.Errorf("invalid ticker: %w", err)
	}

	yesAsk, err := valueobject.NewPrice(resp.YesAsk)
	if err != nil {
		yesAsk, _ = valueobject.NewPrice(0)
	}

	yesBid, err := valueobject.NewPrice(resp.YesBid)
	if err != nil {
		yesBid, _ = valueobject.NewPrice(0)
	}

	noAsk, err := valueobject.NewPrice(resp.NoAsk)
	if err != nil {
		noAsk, _ = valueobject.NewPrice(0)
	}

	noBid, err := valueobject.NewPrice(resp.NoBid)
	if err != nil {
		noBid, _ = valueobject.NewPrice(0)
	}

	lastPrice, err := valueobject.NewPrice(resp.LastPrice)
	if err != nil {
		lastPrice, _ = valueobject.NewPrice(0)
	}

	status := m.mapMarketStatus(resp.Status)

	market := entity.NewMarket(
		ticker,
		resp.Title,
		resp.Category,
		resp.OpenTime,
		resp.CloseTime,
		status,
	)

	market.YesAsk = yesAsk
	market.YesBid = yesBid
	market.NoAsk = noAsk
	market.NoBid = noBid
	market.LastPrice = lastPrice
	market.Volume = resp.Volume
	market.Volume24h = resp.Volume24h
	market.Liquidity = resp.Liquidity

	return market, nil
}

// ToMarketEntities converts multiple MarketResponse to Market entities
func (m *Mapper) ToMarketEntities(responses []MarketResponse) ([]*entity.Market, error) {
	markets := make([]*entity.Market, 0, len(responses))

	for _, resp := range responses {
		market, err := m.ToMarketEntity(&resp)
		if err != nil {
			continue
		}
		markets = append(markets, market)
	}

	return markets, nil
}

// ToOrderBookEntity converts an OrderBookResponse to an OrderBook entity
func (m *Mapper) ToOrderBookEntity(resp *OrderBookResponse) (*entity.OrderBook, error) {
	ticker, err := valueobject.NewTicker(resp.Ticker)
	if err != nil {
		return nil, fmt.Errorf("invalid ticker: %w", err)
	}

	bids := make([]entity.OrderLevel, 0, len(resp.YesOrders))
	for _, level := range resp.YesOrders {
		price, err := valueobject.NewPrice(level.Price)
		if err != nil {
			continue
		}
		bids = append(bids, entity.OrderLevel{
			Price:    price,
			Quantity: int(level.Quantity),
		})
	}

	asks := make([]entity.OrderLevel, 0, len(resp.NoOrders))
	for _, level := range resp.NoOrders {
		price, err := valueobject.NewPrice(level.Price)
		if err != nil {
			continue
		}
		asks = append(asks, entity.OrderLevel{
			Price:    price,
			Quantity: int(level.Quantity),
		})
	}

	orderBook := entity.NewOrderBook(ticker, resp.LastUpdate, bids, asks)
	return orderBook, nil
}

// ToTradeEntities converts multiple TradeResponse to Trade entities
func (m *Mapper) ToTradeEntities(responses []TradeResponse) ([]*entity.Trade, error) {
	trades := make([]*entity.Trade, 0, len(responses))

	for _, resp := range responses {
		ticker, err := valueobject.NewTicker(resp.Ticker)
		if err != nil {
			continue
		}

		price, err := valueobject.NewPrice(resp.Price)
		if err != nil {
			continue
		}

		side := m.mapTradeSide(resp.Side)

		trade := entity.NewTrade(
			ticker,
			resp.TradeID,
			price,
			int(resp.Quantity),
			side,
			resp.CreatedAt,
		)

		trades = append(trades, trade)
	}

	return trades, nil
}

// mapMarketStatus converts API status to domain status
func (m *Mapper) mapMarketStatus(status string) entity.MarketStatus {
	switch status {
	case "open":
		return entity.MarketStatusOpen
	case "closed":
		return entity.MarketStatusClosed
	case "settled":
		return entity.MarketStatusSettled
	default:
		return entity.MarketStatusClosed
	}
}

// mapTradeSide converts API side to domain side
func (m *Mapper) mapTradeSide(side string) entity.TradeSide {
	switch side {
	case "yes", "buy":
		return entity.TradeSideBuy
	case "no", "sell":
		return entity.TradeSideSell
	default:
		return entity.TradeSideBuy
	}
}
