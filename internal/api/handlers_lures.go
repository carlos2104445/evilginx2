package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/pkg/models"
)

func (h *Handlers) listLures(c *gin.Context) {
	lures, err := h.storage.ListLures(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"lures": lures})
}

func (h *Handlers) createLure(c *gin.Context) {
	var lure models.Lure
	if err := c.ShouldBindJSON(&lure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.storage.CreateLure(c.Request.Context(), &lure)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, lure)
}

func (h *Handlers) getLure(c *gin.Context) {
	id := c.Param("id")
	
	lure, err := h.storage.GetLure(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lure not found"})
		return
	}
	
	c.JSON(http.StatusOK, lure)
}

func (h *Handlers) updateLure(c *gin.Context) {
	id := c.Param("id")
	
	var lure models.Lure
	if err := c.ShouldBindJSON(&lure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	lure.ID = id
	err := h.storage.UpdateLure(c.Request.Context(), &lure)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, lure)
}

func (h *Handlers) deleteLure(c *gin.Context) {
	id := c.Param("id")
	
	err := h.storage.DeleteLure(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Lure deleted successfully"})
}
