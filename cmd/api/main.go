package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"honnef.co/go/tools/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting Kalshi Aggregation API on port %s (mode: %s)\n",
		cfg.Server.Port, cfg.Server.GinMode)

	// Initialize Redis client
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

	// Create token service
	tokenService := service.NewTokenService(cfg.JWT.Secret, cfg.JWT.Expiration)
	fmt.Printf("Token service initialized (expiration: %s)\n", cfg.JWT.Expiration.String())

	// Create rate limiter
	rateLimitRepo := ratelimit.NewRedisRateLimiter(redisClient)
	rateLimiter := ratelimitservice.NewRateLimiter(rateLimitRepo)
	fmt.Println("Rate limiter initialized")

	// Create Kalshi API client
	kalshiClient := kalshi.NewClient(cfg.Kalshi.BaseURL, cfg.Kalshi.APIKey)
	fmt.Println("Kalshi API client initialized")

	// Create market repository
	marketRepo := cache.NewMarketRepository(redisClient, kalshiClient)
	fmt.Println("Market repository initialized")

	// Create use cases
	listMarketsUseCase := usecase.NewListMarkets(marketRepo)
	getMarketDetailsUseCase := usecase.NewGetMarketDetails(marketRepo)
	fmt.Println("Use cases initialized")

	// Create HTTP server
	server := httpserver.NewServer(cfg, redisClient, tokenService, rateLimiter, listMarketsUseCase, getMarketDetailsUseCase)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server failed to start: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Received shutdown signal")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server exited gracefully")
}
