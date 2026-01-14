package kalshi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout    = 60 * time.Second // Increased for slow /series endpoint
	maxRetries        = 3
	initialBackoff    = 1 * time.Second
	maxBackoff        = 10 * time.Second
	backoffMultiplier = 2.0
)

// Client represents a Kalshi API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Kalshi API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// GetMarkets fetches markets filtered by category and status
func (c *Client) GetMarkets(ctx context.Context, category string, status string) (*MarketListResponse, error) {
	// Step 1: Get series tickers for this category
	seriesTickers, err := c.getSeriesTickersForCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}
	
	if len(seriesTickers) == 0 {
		return &MarketListResponse{Markets: []MarketResponse{}}, nil
	}
	
	// Step 2: Fetch markets for each series individually
	// Note: Kalshi API doesn't support multiple series_ticker params in one call
	// We'll fetch markets for each series (limited to first 50 series to avoid timeout)
	upperCategory := strings.ToUpper(category)
	
	var allMarkets []MarketResponse
	maxSeries := 10 // Limit to avoid timeout for categories with many series
	if len(seriesTickers) > maxSeries {
		seriesTickers = seriesTickers[:maxSeries]
	}
	
	for _, seriesTicker := range seriesTickers {
		url := fmt.Sprintf("%s/trade-api/v2/markets?series_ticker=%s&limit=1", c.baseURL, seriesTicker)
		if status != "" {
			url += fmt.Sprintf("&status=%s", status)
		}
		
		var response MarketListResponse
		if err := c.doRequest(ctx, "GET", url, nil, &response); err != nil {
			// Log error but continue with other series
			continue
		}
		
		// Add markets from this series
		for i := range response.Markets {
			response.Markets[i].Category = upperCategory
		}
		allMarkets = append(allMarkets, response.Markets...)
	}	
	
	return &MarketListResponse{Markets: allMarkets}, nil
}

// getSeriesTickersForCategory fetches series and returns tickers for the given category
func (c *Client) getSeriesTickersForCategory(ctx context.Context, category string) ([]string, error) {
	// Fetch series with smaller limit to avoid timeout
	// The /series endpoint is slow; limiting to 500 keeps response time reasonable
	url := fmt.Sprintf("%s/trade-api/v2/series?category=%s", c.baseURL, category)
	
	var response SeriesListResponse
	if err := c.doRequest(ctx, "GET", url, nil, &response); err != nil {
		return nil, err
	}
	
	upperCategory := strings.ToUpper(category)
	seriesTickers := make([]string, 0)
	
	for _, series := range response.Series {
		if strings.EqualFold(series.Category, upperCategory) {
			seriesTickers = append(seriesTickers, series.Ticker)
		}
	}
	
	return seriesTickers, nil
}

// GetMarket fetches a single market by tickery
func (c *Client) GetMarket(ctx context.Context, ticker string) (*MarketResponse, error) {
	url := fmt.Sprintf("%s/trade-api/v2/markets/%s", c.baseURL, ticker)

	var response struct {
		Market MarketResponse `json:"market"`
	}
	if err := c.doRequest(ctx, "GET", url, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get market: %w", err)
	}

	return &response.Market, nil
}

// GetOrderBook fetches the order book for a market
func (c *Client) GetOrderBook(ctx context.Context, ticker string) (*OrderBookResponse, error) {
	url := fmt.Sprintf("%s/trade-api/v2/markets/%s/orderbook", c.baseURL, ticker)

	var response OrderBookResponse
	if err := c.doRequest(ctx, "GET", url, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get orderbook: %w", err)
	}

	return &response, nil
}

// GetTrades fetches recent trades for a market
func (c *Client) GetTrades(ctx context.Context, ticker string, limit int) (*TradesResponse, error) {
	url := fmt.Sprintf("%s/trade-api/v2/markets/%s/trades?limit=%d", c.baseURL, ticker, limit)

	var response TradesResponse
	if err := c.doRequest(ctx, "GET", url, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}

	return &response, nil
}

// doRequest executes an HTTP request with retry logic and exponential backoff
func (c *Client) doRequest(ctx context.Context, method, url string, body io.Reader, result interface{}) error {
	var lastErr error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				// Calculate next backoff
				backoff = time.Duration(float64(backoff) * backoffMultiplier)
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		
		// Only set Authorization header if API key is provided
		if c.apiKey != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			continue
		}

		if resp.StatusCode == 429 {
			// Rate limit - retry with backoff
			lastErr = fmt.Errorf("rate limit exceeded")
			continue
		}

		if resp.StatusCode >= 400 {
			// Client error - don't retry
			var errResp ErrorResponse
			if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Error.Message != "" {
				return fmt.Errorf("API error: %s (code: %s)", errResp.Error.Message, errResp.Error.Code)
			}
			return fmt.Errorf("API error: status %d", resp.StatusCode)
		}

		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return nil
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}
