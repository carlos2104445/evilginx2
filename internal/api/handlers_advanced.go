package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/internal/phishlet"
	"github.com/kgretzky/evilginx2/pkg/models"
)

type EvaluateConditionsRequest struct {
	UserAgent    string            `json:"user_agent"`
	Email        string            `json:"email"`
	IPAddress    string            `json:"ip_address"`
	Hostname     string            `json:"hostname"`
	Path         string            `json:"path"`
	CustomParams map[string]string `json:"custom_params"`
}

type CreateVersionRequest struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

type UpdateFlowStepRequest struct {
	SessionID string            `json:"session_id"`
	StepData  map[string]string `json:"step_data"`
}

func (h *Handlers) listPhishletVersions(c *gin.Context) {
	name := c.Param("name")
	
	versions, err := h.phishletRepo.ListVersions(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

func (h *Handlers) createPhishletVersion(c *gin.Context) {
	name := c.Param("name")
	
	var req CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	phishlet, err := h.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phishlet not found"})
		return
	}
	
	err = h.phishletRepo.PublishVersion(c.Request.Context(), phishlet, req.Version, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Version created successfully"})
}

func (h *Handlers) getPhishletVersion(c *gin.Context) {
	name := c.Param("name")
	version := c.Param("version")
	
	phishlet, err := h.phishletRepo.GetVersion(c.Request.Context(), name, version)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, phishlet)
}

func (h *Handlers) evaluateConditions(c *gin.Context) {
	name := c.Param("name")
	
	var req EvaluateConditionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	phishlet, err := h.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phishlet not found"})
		return
	}
	
	evaluator := phishlet.NewConditionEvaluator()
	context := &phishlet.EvaluationContext{
		UserAgent:    req.UserAgent,
		Email:        req.Email,
		IPAddress:    req.IPAddress,
		Hostname:     req.Hostname,
		Path:         req.Path,
		CustomParams: req.CustomParams,
	}
	
	actions, err := evaluator.EvaluateConditions(phishlet, context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"actions": actions})
}

func (h *Handlers) getMultiPageFlows(c *gin.Context) {
	name := c.Param("name")
	
	phishlet, err := h.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phishlet not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"flows": phishlet.MultiPageFlows})
}

func (h *Handlers) updateFlowStep(c *gin.Context) {
	name := c.Param("name")
	flowName := c.Param("flow")
	
	var req UpdateFlowStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.phishletRepo.UpdateFlowSession(c.Request.Context(), req.SessionID, flowName, req.StepData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Flow step updated successfully"})
}
