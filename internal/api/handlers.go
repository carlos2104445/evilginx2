package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/internal/utils"
)

func (h *Handlers) createSession(c *gin.Context) {
	var req struct {
		PhishletName string `json:"phishlet_name" binding:"required"`
		UserAgent    string `json:"user_agent"`
		RemoteAddr   string `json:"remote_addr"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.PhishletName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phishlet name cannot be empty"})
		return
	}

	if len(req.PhishletName) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phishlet name too long"})
		return
	}

	req.PhishletName = utils.SanitizeString(req.PhishletName, 100)
	req.UserAgent = utils.SanitizeString(req.UserAgent, 500)

	if req.RemoteAddr != "" {
		if err := utils.ValidateIPAddress(req.RemoteAddr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address format"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Session creation validated",
		"data": gin.H{
			"phishlet_name": req.PhishletName,
			"user_agent":    req.UserAgent,
			"remote_addr":   req.RemoteAddr,
		},
	})
}

func (h *Handlers) listSessions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 || limit > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	phishletFilter := c.Query("phishlet")
	if phishletFilter != "" {
		phishletFilter = utils.SanitizeString(phishletFilter, 100)
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": []gin.H{},
		"total":    0,
		"limit":    limit,
		"offset":   offset,
		"filters": gin.H{
			"phishlet": phishletFilter,
		},
	})
}

func (h *Handlers) getSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID cannot be empty"})
		return
	}

	sessionID = utils.SanitizeString(sessionID, 64)

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"status":     "active",
	})
}

func (h *Handlers) createPhishlet(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Author      string `json:"author"`
		Version     string `json:"version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phishlet name cannot be empty"})
		return
	}

	if len(req.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phishlet name too long"})
		return
	}

	if strings.Contains(req.Name, "/") || strings.Contains(req.Name, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid characters in phishlet name"})
		return
	}

	req.Name = utils.SanitizeString(req.Name, 100)
	req.Description = utils.SanitizeString(req.Description, 500)
	req.Author = utils.SanitizeString(req.Author, 100)
	req.Version = utils.SanitizeString(req.Version, 20)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Phishlet creation validated",
		"data": gin.H{
			"name":        req.Name,
			"description": req.Description,
			"author":      req.Author,
			"version":     req.Version,
		},
	})
}

func (h *Handlers) listPhishlets(c *gin.Context) {
	enabledOnly := c.Query("enabled") == "true"
	
	c.JSON(http.StatusOK, gin.H{
		"phishlets": []gin.H{},
		"total":     0,
		"filters": gin.H{
			"enabled_only": enabledOnly,
		},
	})
}

func (h *Handlers) getPhishlet(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phishlet name cannot be empty"})
		return
	}

	name = utils.SanitizeString(name, 100)

	c.JSON(http.StatusOK, gin.H{
		"name":    name,
		"enabled": false,
		"status":  "inactive",
	})
}
