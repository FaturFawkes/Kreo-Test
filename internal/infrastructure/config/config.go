package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server    ServerConfig
	Redis     RedisConfig
	Kalshi    KalshiConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Cache     CacheConfig
	Worker    WorkerConfig
	Logging   LoggingConfig
}

type ServerConfig struct {
	Port         string
	GinMode      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type KalshiConfig struct {
	BaseURL string
	APIKey  string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type RateLimitConfig struct {
	Authenticated   int
	Unauthenticated int
	Worker          int
}

type CacheConfig struct {
	TTLMarkets  time.Duration
	TTLDetails  time.Duration
	TTLOverview time.Duration
}

type WorkerConfig struct {
	PoolSize        int
	IntervalSeconds int
	HotMarketCount  int
}

type LoggingConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			GinMode:      getEnv("GIN_MODE", "release"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Kalshi: KalshiConfig{
			BaseURL: getEnv("KALSHI_API_BASE_URL", "https://api.kalshi.com"),
			APIKey:  getEnv("KALSHI_API_KEY", ""),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Expiration: time.Duration(getEnvInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		},
		RateLimit: RateLimitConfig{
			Authenticated:   getEnvInt("RATE_LIMIT_AUTHENTICATED", 100),
			Unauthenticated: getEnvInt("RATE_LIMIT_UNAUTHENTICATED", 10),
			Worker:          getEnvInt("RATE_LIMIT_WORKER", 80),
		},
		Cache: CacheConfig{
			TTLMarkets:  time.Duration(getEnvInt("CACHE_TTL_MARKETS", 300)) * time.Second,
			TTLDetails:  time.Duration(getEnvInt("CACHE_TTL_DETAILS", 60)) * time.Second,
			TTLOverview: time.Duration(getEnvInt("CACHE_TTL_OVERVIEW", 300)) * time.Second,
		},
		Worker: WorkerConfig{
			PoolSize:        getEnvInt("WORKER_POOL_SIZE", 5),
			IntervalSeconds: getEnvInt("WORKER_INTERVAL_SECONDS", 60),
			HotMarketCount:  getEnvInt("HOT_MARKET_COUNT", 20),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validate required fields
	if cfg.Kalshi.APIKey == "" {
		return nil, fmt.Errorf("KALSHI_API_KEY is required")
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as an integer or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// RedisAddr returns the Redis connection address
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
