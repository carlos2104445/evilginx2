package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/internal/storage"
	"github.com/kgretzky/evilginx2/pkg/models"
)

func (s *Server) listPhishlets(c *gin.Context) {
	filters := &storage.PhishletFilters{}
	
	if name := c.Query("name"); name != "" {
		filters.Name = name
	}
	
	if enabled := c.Query("enabled"); enabled != "" {
		if e, err := strconv.ParseBool(enabled); err == nil {
			filters.Enabled = &e
		}
	}
	
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filters.Limit = l
		}
	}
	
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filters.Offset = o
		}
	}
	
	phishlets, err := s.storage.ListPhishlets(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"phishlets": phishlets,
		"count":     len(phishlets),
	})
}

func (s *Server) createPhishlet(c *gin.Context) {
	var phishlet models.Phishlet
	if err := c.ShouldBindJSON(&phishlet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if phishlet.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phishlet name is required"})
		return
	}
	
	if err := s.storage.CreatePhishlet(c.Request.Context(), &phishlet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, phishlet)
}

func (s *Server) getPhishlet(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phishlet name is required"})
		return
	}
	
	phishlet, err := s.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "phishlet not found"})
		return
	}
	
	c.JSON(http.StatusOK, phishlet)
}

func (s *Server) updatePhishlet(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phishlet name is required"})
		return
	}
	
	var phishlet models.Phishlet
	if err := c.ShouldBindJSON(&phishlet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	phishlet.Name = name
	
	if err := s.storage.UpdatePhishlet(c.Request.Context(), &phishlet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, phishlet)
}

func (s *Server) deletePhishlet(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phishlet name is required"})
		return
	}
	
	if err := s.storage.DeletePhishlet(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "phishlet deleted successfully"})
}

func (s *Server) getPhishletStats(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phishlet name is required"})
		return
	}
	
	sessionFilters := &storage.SessionFilters{
		PhishletName: name,
	}
	
	sessions, err := s.storage.ListSessions(c.Request.Context(), sessionFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	stats := models.PhishletStats{
		TotalPhishlets:   1,
		EnabledPhishlets: 1,
		ActiveCampaigns:  len(sessions),
	}
	
	c.JSON(http.StatusOK, stats)
}
