package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client with the given configuration
func NewRedisClient(addr string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,

		// Connection pool settings
		PoolSize:     10,
		MinIdleConns: 5,

		// Timeouts - increased toa handle slow Kalshi API calls
		DialTimeout:  5 * time.Second,
		ReadTimeout:  90 * time.Second, // Increased from 3s to 90s
		WriteTimeout: 90 * time.Second, // Increased from 3s to 90s
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
