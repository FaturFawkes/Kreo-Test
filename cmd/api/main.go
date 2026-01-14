package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"upwork-test/internal/application/usecase"
	httpserver "upwork-test/internal/delivery/http"
	"upwork-test/internal/domain/auth/service"
	ratelimitservice "upwork-test/internal/domain/ratelimit/service"
	"upwork-test/internal/infrastructure/cache"
	"upwork-test/internal/infrastructure/config"
	"upwork-test/internal/infrastructure/kalshi"
	"upwork-test/internal/infrastructure/ratelimit"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting Kalshi Aggregation API on port %s (mode: %s)\n",
		cfg.Server.Port, cfg.Server.GinMode)

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

	tokenService := service.NewTokenService(cfg.JWT.Secret, cfg.JWT.Expiration)
	fmt.Printf("Token service initialized (expiration: %s)\n", cfg.JWT.Expiration.String())

	rateLimitRepo := ratelimit.NewRedisRateLimiter(redisClient)
	rateLimiter := ratelimitservice.NewRateLimiter(rateLimitRepo)
	fmt.Println("Rate limiter initialized")

	kalshiClient := kalshi.NewClient(cfg.Kalshi.BaseURL, cfg.Kalshi.APIKey)
	fmt.Println("Kalshi API client initialized")

	marketRepo := cache.NewMarketRepository(redisClient, kalshiClient)
	fmt.Println("Market repository initialized")

	categoryRepo := cache.NewCategoryRepository(redisClient, kalshiClient, marketRepo)
	fmt.Println("Category repository initialized")

	listMarketsUseCase := usecase.NewListMarkets(marketRepo)
	getMarketDetailsUseCase := usecase.NewGetMarketDetails(marketRepo)
	getCategoryOverviewUseCase := usecase.NewGetCategoryOverview(categoryRepo)
	fmt.Println("Use cases initialized")

	server := httpserver.NewServer(cfg, redisClient, tokenService, rateLimiter, listMarketsUseCase, getMarketDetailsUseCase, getCategoryOverviewUseCase)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server failed to start: %v\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server exited gracefully")
}
