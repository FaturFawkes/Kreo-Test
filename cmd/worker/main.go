package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"upwork-test/internal/application/service"
	"upwork-test/internal/infrastructure/cache"
	"upwork-test/internal/infrastructure/config"
	"upwork-test/internal/infrastructure/kalshi"

	"golang.org/x/time/rate"
)

const (
	maxWorkers      = 5
	kalshiRateLimit = 80 // requests per minute (80% of assumed 100 req/min limit)
	hotMarketsCount = 20
)

func main() {
	fmt.Println("Starting cache warmer worker...")

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	redisClient, err := cache.NewRedisClient(
		cfg.Redis.Addr(),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	fmt.Printf("Connected to Redis at %s\n", cfg.Redis.Addr())

	kalshiClient := kalshi.NewClient(cfg.Kalshi.BaseURL, cfg.Kalshi.APIKey)

	marketRepo := cache.NewMarketRepository(redisClient, kalshiClient)
	categoryRepo := cache.NewCategoryRepository(redisClient, kalshiClient, marketRepo)

	cacheWarmer := service.NewCacheWarmer(marketRepo, categoryRepo)

	limiter := rate.NewLimiter(rate.Every(time.Minute/kalshiRateLimit), kalshiRateLimit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()

		fmt.Println("Initial hot markets warm-up...")
		if err := limiter.Wait(ctx); err == nil {
			if err := cacheWarmer.WarmHotMarkets(ctx, hotMarketsCount); err != nil {
				fmt.Printf("Error warming hot markets: %v\n", err)
			} else {
				fmt.Println("Hot markets warmed successfully")
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Printf("[%s] Warming hot markets...\n", time.Now().Format(time.RFC3339))
				if err := limiter.Wait(ctx); err == nil {
					if err := cacheWarmer.WarmHotMarkets(ctx, hotMarketsCount); err != nil {
						fmt.Printf("Error warming hot markets: %v\n", err)
					}
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		fmt.Println("Initial category overviews warm-up...")
		if err := limiter.Wait(ctx); err == nil {
			if err := cacheWarmer.WarmCategoryOverviews(ctx); err != nil {
				fmt.Printf("Error warming category overviews: %v\n", err)
			} else {
				fmt.Println("Category overviews warmed successfully")
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Printf("[%s] Warming category overviews...\n", time.Now().Format(time.RFC3339))
				if err := limiter.Wait(ctx); err == nil {
					if err := cacheWarmer.WarmCategoryOverviews(ctx); err != nil {
						fmt.Printf("Error warming category overviews: %v\n", err)
					}
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		fmt.Println("Initial category list warm-up...")
		if err := limiter.Wait(ctx); err == nil {
			if err := cacheWarmer.WarmCategoryLists(ctx); err != nil {
				fmt.Printf("Error warming category lists: %v\n", err)
			} else {
				fmt.Println("Category lists warmed successfully")
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Printf("[%s] Warming category lists...\n", time.Now().Format(time.RFC3339))
				if err := limiter.Wait(ctx); err == nil {
					if err := cacheWarmer.WarmCategoryLists(ctx); err != nil {
						fmt.Printf("Error warming category lists: %v\n", err)
					}
				}
			}
		}
	}()

	fmt.Println("Cache warmer workers started successfully")

	<-quit
	fmt.Println("Shutting down worker...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All workers stopped gracefully")
	case <-shutdownCtx.Done():
		fmt.Println("Shutdown timeout exceeded, forcing exit")
	}

	fmt.Println("Worker exited")
}
