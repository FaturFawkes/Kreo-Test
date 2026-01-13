package http

import (
	"context"
	"fmt"
	"net/http"
	"upwork-test/internal/delivery/http/middleware"
	"upwork-test/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
	httpServer              *http.Server
}

// NewServer creates a new HTTP server
func NewServer(
	cfg *config.Config,
) *Server {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)
	router := gin.New()

	// Create server
	srv := &Server{
		config: cfg,
		router: router,
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
	s.router.Use(middleware.ErrorHandler())
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Public routes

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
