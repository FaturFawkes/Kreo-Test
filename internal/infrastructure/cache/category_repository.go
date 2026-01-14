package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"upwork-test/internal/domain/category/entity"
	"upwork-test/internal/domain/category/repository"
	"upwork-test/internal/domain/category/valueobject"
	marketrepo "upwork-test/internal/domain/market/repository"
	"upwork-test/internal/infrastructure/kalshi"

	"github.com/redis/go-redis/v9"
)

const (
	categoryListCacheTTL     = 24 * time.Hour
	categoryOverviewCacheTTL = 10 * time.Minute
)

// CategoryRepository implements the category repository with Redis caching.
type CategoryRepository struct {
	redisClient  *redis.Client
	kalshiClient *kalshi.Client
	marketRepo   marketrepo.MarketRepository
	keyBuilder   *KeyBuilder
}

// NewCategoryRepository creates a new category repository.
func NewCategoryRepository(redisClient *redis.Client, kalshiClient *kalshi.Client, marketRepo marketrepo.MarketRepository) *CategoryRepository {
	return &CategoryRepository{
		redisClient:  redisClient,
		kalshiClient: kalshiClient,
		marketRepo:   marketRepo,
		keyBuilder:   NewKeyBuilder("kalshi"),
	}
}

// GetAll retrieves all available categories.
func (r *CategoryRepository) GetAll(ctx context.Context) ([]*entity.Category, error) {
	cacheKey := r.keyBuilder.CategoryList()

	cachedData, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var categories []*entity.Category
		if err := json.Unmarshal([]byte(cachedData), &categories); err == nil {
			return categories, nil
		}
	}

	categories := r.buildCategoryList()

	if data, err := json.Marshal(categories); err == nil {
		r.redisClient.Set(ctx, cacheKey, data, categoryListCacheTTL)
	}

	return categories, nil
}

func (r *CategoryRepository) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	categoryName, err := valueobject.NewCategoryName(name)
	if err != nil {
		return nil, repository.ErrCategoryNotFound
	}

	categories, err := r.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	for _, cat := range categories {
		if cat.Name.Equals(categoryName) {
			return cat, nil
		}
	}

	return nil, repository.ErrCategoryNotFound
}

func (r *CategoryRepository) GetOverview(ctx context.Context, categoryName string) (*entity.CategoryOverview, error) {
	cacheKey := r.keyBuilder.CategoryOverview(categoryName)

	cachedData, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var overview entity.CategoryOverview
		if err := json.Unmarshal([]byte(cachedData), &overview); err == nil {
			if overview.IsFresh() {
				return &overview, nil
			}
		}
	}

	overview, err := r.computeOverview(ctx, categoryName)
	if err != nil {
		return nil, err
	}

	if err := r.SaveOverview(ctx, overview); err != nil {
		fmt.Printf("Warning: failed to cache overview: %v\n", err)
	}

	return overview, nil
}

func (r *CategoryRepository) SaveOverview(ctx context.Context, overview *entity.CategoryOverview) error {
	cacheKey := r.keyBuilder.CategoryOverview(overview.CategoryName.String())

	data, err := json.Marshal(overview)
	if err != nil {
		return fmt.Errorf("failed to marshal overview: %w", err)
	}

	if err := r.redisClient.Set(ctx, cacheKey, data, categoryOverviewCacheTTL).Err(); err != nil {
		return fmt.Errorf("failed to cache overview: %w", err)
	}

	return nil
}

func (r *CategoryRepository) computeOverview(ctx context.Context, categoryName string) (*entity.CategoryOverview, error) {
	catName, err := valueobject.NewCategoryName(categoryName)
	if err != nil {
		return nil, repository.ErrCategoryNotFound
	}

	markets, total, err := r.marketRepo.ListByCategory(ctx, categoryName, 1, 1000, "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch markets: %w", err)
	}

	var totalVolume24h int64
	var totalLiquidity int64
	activeMarkets := 0

	for _, market := range markets {
		totalVolume24h += market.Volume24h
		totalLiquidity += market.Liquidity
		if market.IsOpen() {
			activeMarkets++
		}
	}

	var avgLiquidity float64
	if len(markets) > 0 {
		avgLiquidity = float64(totalLiquidity) / float64(len(markets))
	}

	// ActiveTraders24h not available from Kalshi API, would come from analytics service in production
	overview, err := entity.NewCategoryOverview(
		catName,
		total,
		totalVolume24h,
		avgLiquidity,
		0, // ActiveTraders24h - not available from current API
		categoryOverviewCacheTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create overview: %w", err)
	}

	return overview, nil
}

func (r *CategoryRepository) buildCategoryList() []*entity.Category {
	categoryNames := valueobject.AllCategories()
	categories := make([]*entity.Category, 0, len(categoryNames))

	categoryDisplayNames := map[string]string{
		"ECONOMICS":     "Economics",
		"POLITICS":      "Politics",
		"SPORTS":        "Sports",
		"CLIMATE":       "Climate",
		"SCIENCE":       "Science",
		"TECHNOLOGY":    "Technology",
		"ENTERTAINMENT": "Entertainment",
		"FINANCE":       "Finance",
		"HEALTH":        "Health",
		"CRYPTO":        "Crypto",
	}

	categoryDescriptions := map[string]string{
		"ECONOMICS":     "Economic indicators and market predictions",
		"POLITICS":      "Political events and election outcomes",
		"SPORTS":        "Sports events and championship predictions",
		"CLIMATE":       "Climate and environmental event predictions",
		"SCIENCE":       "Scientific discoveries and research outcomes",
		"TECHNOLOGY":    "Technology trends and product launches",
		"ENTERTAINMENT": "Entertainment industry events and awards",
		"FINANCE":       "Financial markets and company performance",
		"HEALTH":        "Public health outcomes and medical advances",
		"CRYPTO":        "Cryptocurrency and blockchain predictions",
	}

	for _, name := range categoryNames {
		catName, _ := valueobject.NewCategoryName(name)
		displayName := categoryDisplayNames[name]
		if displayName == "" {
			displayName = name
		}
		description := categoryDescriptions[name]

		category := entity.NewCategory(catName, displayName, description)
		categories = append(categories, category)
	}

	return categories
}
