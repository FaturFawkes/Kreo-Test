# Kalshi Aggregation API

High-performance, production-ready Golang API that exposes aggregated Kalshi prediction market data using the Kalshi API as a read-only source of truth.

## Overview

This API provides cached access to Kalshi prediction market data with the following features:

- **High Performance**: 50ms p95 latency for cached responses, 90% cache hit rate
- **Production Ready**: JWT authentication, rate limiting, structured logging, health checks
- **Cache-First Architecture**: Redis-based caching with intelligent cache warming
- **Domain-Driven Design**: Clean architecture with bounded contexts for maintainability

## Architecture

The system follows Domain-Driven Design with these bounded contexts:

- **Market**: Market listing, details aggregation, order book, trade history
- **Category**: Category-based market browsing, overview metrics
- **Auth**: JWT-based authentication with token management
- **RateLimit**: Request rate limiting with tiered limits

### Technology Stack

- **Language**: Go 1.25+
- **Framework**: Gin for HTTP routing
- **Cache**: Redis 7+ for caching and rate limiting
- **Authentication**: JWT with 24-hour token expiration
- **Logging**: logmanager with structured logging and trace IDs
- **Testing**: testify for assertions, mockery for mocks
- **Containerization**: Docker with multi-stage builds

## API Endpoints

### Authentication
- `POST /auth/token` - Generate JWT token (requires API credentials)

### Markets
- `GET /categories/{category}/markets` - List markets in a category
- `GET /markets/{ticker}` - Get aggregated market details (metadata + orderbook + trades)

### Categories
- `GET /categories/{category}/overview` - Get category overview metrics

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.25+ (for local development)
- Redis 7+ (managed by Docker Compose)

### Running with Docker Compose

1. Create `.env` file:
```bash
cat > .env << EOF
# Kalshi API Configuration
KALSHI_API_BASE_URL=https://api.kalshi.com
KALSHI_API_KEY=your-api-key-here

# JWT Configuration
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION_HOURS=24

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Server Configuration
PORT=8080
GIN_MODE=release

# Rate Limiting
RATE_LIMIT_AUTHENTICATED=100
RATE_LIMIT_UNAUTHENTICATED=10
RATE_LIMIT_WORKER=80

# Cache Configuration
CACHE_TTL_MARKETS=300
CACHE_TTL_DETAILS=60
CACHE_TTL_OVERVIEW=300
EOF
```

2. Start services:
```bash
docker-compose up -d
```

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Run Redis:
```bash
docker run -d -p 6379:6379 redis:7-alpine
```

3. Set environment variables (use `.env` file or export):
```bash
export REDIS_HOST=localhost
export REDIS_PORT=6379
export PORT=8080
# ... other variables from .env
```

4. Run API server:
```bash
go run cmd/api/main.go
```

5. Run worker (in separate terminal):
```bash
go run cmd/worker/main.go
```

## Development

### Project Structure

```
.
├── cmd/
│   ├── api/              # API server entry point
│   └── worker/           # Background worker entry point
├── internal/
│   ├── domain/           # Domain layer (entities, value objects, repositories)
│   │   ├── market/
│   │   ├── category/
│   │   ├── auth/
│   │   └── ratelimit/
│   ├── application/      # Application layer (use cases, services, DTOs)
│   ├── infrastructure/   # Infrastructure layer (Kalshi client, Redis, config)
│   └── delivery/         # Delivery layer (HTTP handlers, middleware)
├── specs/                # Feature specifications and documentation
├── .golangci.yaml        # Linter configuration
├── .mockery.yaml         # Mock generation configuration
├── Dockerfile            # Multi-stage container build
└── docker-compose.yaml   # Service orchestration
```

## Performance Goals

- **Latency**: 50ms p95 for cached responses, 500ms for cache misses
- **Throughput**: 1000 concurrent requests
- **Cache Hit Rate**: 90%
- **Availability**: 99.9% uptime

## Rate Limits

- **Authenticated Users**: 100 requests/minute
- **Unauthenticated Users**: 10 requests/minute  
- **Background Workers**: 80 requests/minute to Kalshi API

## License

Proprietary - All rights reserved

## Support

For issues or questions, please contact the development team.
