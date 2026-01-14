package service

import (
	"context"
	"fmt"
	"sort"
	"time"
	"upwork-test/internal/domain/category/repository"
	marketrepo "upwork-test/internal/domain/market/repository"
)

// CacheWarmer handles proactive cache refreshing
type CacheWarmer struct {
	marketRepo   marketrepo.MarketRepository
	categoryRepo repository.CategoryRepository
}

// NewCacheWarmer creates a new cache warmer
func NewCacheWarmer(
	marketRepo marketrepo.MarketRepository,
	categoryRepo repository.CategoryRepository,
) *CacheWarmer {
	return &CacheWarmer{
		marketRepo:   marketRepo,
		categoryRepo: categoryRepo,
	}
}

// WarmHotMarkets refreshes cache for top markets by volume
func (cw *CacheWarmer) WarmHotMarkets(ctx context.Context, topN int) error {
	categories, err := cw.categoryRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get categories: %w", err)
	}

	type marketWithVolume struct {
		ticker    string
		volume24h int64
	}
	var allMarkets []marketWithVolume

	for _, cat := range categories {
		markets, _, err := cw.marketRepo.ListByCategory(ctx, cat.Name.String(), 1, 200, "")
		if err != nil {
			fmt.Printf("Warning: failed to get markets for category %s: %v\n", cat.Name.String(), err)
			continue
		}

		for _, market := range markets {
			allMarkets = append(allMarkets, marketWithVolume{
				ticker:    market.Ticker.String(),
				volume24h: market.Volume24h,
			})
		}
	}

	sort.Slice(allMarkets, func(i, j int) bool {
		return allMarkets[i].volume24h > allMarkets[j].volume24h
	})

	limit := topN
	if len(allMarkets) < limit {
		limit = len(allMarkets)
	}

	for i := 0; i < limit; i++ {
		ticker := allMarkets[i].ticker

		_, err := cw.marketRepo.GetByTicker(ctx, ticker)
		if err != nil {
			fmt.Printf("Warning: failed to warm market %s: %v\n", ticker, err)
		}

		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// WarmCategoryOverviews refreshes cache for all category overviews
func (cw *CacheWarmer) WarmCategoryOverviews(ctx context.Context) error {
	categories, err := cw.categoryRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get categories: %w", err)
	}

	for _, cat := range categories {
		_, err := cw.categoryRepo.GetOverview(ctx, cat.Name.String())
		if err != nil {
			fmt.Printf("Warning: failed to warm category overview %s: %v\n", cat.Name.String(), err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// WarmCategoryLists refreshes the category list cache
func (cw *CacheWarmer) WarmCategoryLists(ctx context.Context) error {
	_, err := cw.categoryRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to warm category list: %w", err)
	}
	return nil
}

// WarmMarketsByCategory refreshes market list cache for a specific category
func (cw *CacheWarmer) WarmMarketsByCategory(ctx context.Context, category string) error {
	_, _, err := cw.marketRepo.ListByCategory(ctx, category, 1, 200, "")
	if err != nil {
		return fmt.Errorf("failed to warm markets for category %s: %w", category, err)
	}
	return nil
}
