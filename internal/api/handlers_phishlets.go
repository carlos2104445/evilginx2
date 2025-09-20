package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/pkg/models"
)

func (h *Handlers) listPhishlets(c *gin.Context) {
	phishlets, err := h.storage.ListPhishlets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"phishlets": phishlets})
}

func (h *Handlers) createPhishlet(c *gin.Context) {
	var phishlet models.Phishlet
	if err := c.ShouldBindJSON(&phishlet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.storage.CreatePhishlet(c.Request.Context(), &phishlet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, phishlet)
}

func (h *Handlers) getPhishlet(c *gin.Context) {
	name := c.Param("name")
	
	phishlet, err := h.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phishlet not found"})
		return
	}
	
	c.JSON(http.StatusOK, phishlet)
}

func (h *Handlers) updatePhishlet(c *gin.Context) {
	name := c.Param("name")
	
	var phishlet models.Phishlet
	if err := c.ShouldBindJSON(&phishlet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	phishlet.Name = name
	err := h.storage.UpdatePhishlet(c.Request.Context(), &phishlet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, phishlet)
}

func (h *Handlers) deletePhishlet(c *gin.Context) {
	name := c.Param("name")
	
	err := h.storage.DeletePhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Phishlet deleted successfully"})
}
