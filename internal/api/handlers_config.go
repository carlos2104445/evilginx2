package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *Handlers) getConfig(c *gin.Context) {
	config, err := h.storage.ListConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"config": config})
}

func (h *Handlers) setConfig(c *gin.Context) {
	var req ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.storage.SetConfig(c.Request.Context(), req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

func (h *Handlers) deleteConfig(c *gin.Context) {
	key := c.Param("key")
	
	err := h.storage.DeleteConfig(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}
