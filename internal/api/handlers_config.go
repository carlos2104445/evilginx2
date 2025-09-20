package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/pkg/models"
)

func (s *Server) getConfig(c *gin.Context) {
	config, err := s.storage.GetConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, config)
}

func (s *Server) updateConfig(c *gin.Context) {
	var config models.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := s.storage.UpdateConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	s.config = &config
	
	c.JSON(http.StatusOK, config)
}

func (s *Server) listLures(c *gin.Context) {
	lures, err := s.storage.ListLures(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"lures": lures,
		"count": len(lures),
	})
}

func (s *Server) createLure(c *gin.Context) {
	var lure models.Lure
	if err := c.ShouldBindJSON(&lure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if lure.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lure ID is required"})
		return
	}
	
	if err := s.storage.CreateLure(c.Request.Context(), &lure); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, lure)
}

func (s *Server) getLure(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lure ID is required"})
		return
	}
	
	lure, err := s.storage.GetLure(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lure not found"})
		return
	}
	
	c.JSON(http.StatusOK, lure)
}

func (s *Server) updateLure(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lure ID is required"})
		return
	}
	
	var lure models.Lure
	if err := c.ShouldBindJSON(&lure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	lure.ID = id
	
	if err := s.storage.UpdateLure(c.Request.Context(), &lure); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, lure)
}

func (s *Server) deleteLure(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lure ID is required"})
		return
	}
	
	if err := s.storage.DeleteLure(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "lure deleted successfully"})
}
