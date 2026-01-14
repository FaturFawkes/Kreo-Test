package http

import (
	"context"
	"fmt"
	"net/http"
	"upwork-test/internal/application/usecase"
	"upwork-test/internal/delivery/http/handler"
	"upwork-test/internal/delivery/http/middleware"
	"upwork-test/internal/domain/auth/service"
	ratelimitservice "upwork-test/internal/domain/ratelimit/service"
	"upwork-test/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Server represents the HTTP server
type Server struct {
	config                     *config.Config
	router                     *gin.Engine
	httpServer                 *http.Server
	redisClient                *redis.Client
	tokenService               *service.TokenService
	rateLimiter                *ratelimitservice.RateLimiter
	listMarketsUseCase         *usecase.ListMarkets
	getMarketDetailsUseCase    *usecase.GetMarketDetails
	getCategoryOverviewUseCase *usecase.GetCategoryOverview
}

// NewServer creates a new HTTP server
func NewServer(
	cfg *config.Config,
	redisClient *redis.Client,
	tokenService *service.TokenService,
	rateLimiter *ratelimitservice.RateLimiter,
	listMarketsUseCase *usecase.ListMarkets,
	getMarketDetailsUseCase *usecase.GetMarketDetails,
	getCategoryOverviewUseCase *usecase.GetCategoryOverview,
) *Server {
	gin.SetMode(cfg.Server.GinMode)
	router := gin.New()

	srv := &Server{
		config:                     cfg,
		router:                     router,
		redisClient:                redisClient,
		tokenService:               tokenService,
		rateLimiter:                rateLimiter,
		listMarketsUseCase:         listMarketsUseCase,
		getMarketDetailsUseCase:    getMarketDetailsUseCase,
		getCategoryOverviewUseCase: getCategoryOverviewUseCase,
	}

	// Setup middleware and routes
	srv.setupMiddleware()
	srv.setupRoutes()

	// Create HTTP server
	srv.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return srv
}

// setupMiddleware configures the middleware chain
func (s *Server) setupMiddleware() {
	s.router.Use(gin.Recovery())
	s.router.Use(middleware.Logging())
	s.router.Use(middleware.RateLimitMiddleware(s.rateLimiter, s.tokenService))
	s.router.Use(middleware.ErrorHandler())

}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	v1 := s.router.Group("/api/v1")
	{
		// Auth endpoints (public)
		auth := v1.Group("/auth")
		{
			authenticateUseCase := usecase.NewAuthenticate(s.tokenService, s.config.Kalshi.APIKey)
			authHandler := handler.NewAuthHandler(authenticateUseCase)
			auth.POST("/token", authHandler.IssueToken)
		}

		// Protected market endpoints (require authentication)
		categories := v1.Group("/categories")
		categories.Use(middleware.Auth(s.tokenService))
		{
			marketHandler := handler.NewMarketHandler(s.listMarketsUseCase, s.getMarketDetailsUseCase)
			categories.GET("/:category/markets", marketHandler.ListMarkets)

			categoryHandler := handler.NewCategoryHandler(s.getCategoryOverviewUseCase)
			categories.GET("/:category/overview", categoryHandler.GetOverview)
		}

		// Protected market detail endpoint
		markets := v1.Group("/markets")
		markets.Use(middleware.Auth(s.tokenService))
		{
			marketHandler := handler.NewMarketHandler(s.listMarketsUseCase, s.getMarketDetailsUseCase)
			markets.GET("/:ticker", marketHandler.GetMarketDetails)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("Starting HTTP server on port %s\n", s.config.Server.Port)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
