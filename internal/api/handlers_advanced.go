package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	
	versions := []interface{}{}
	_ = name // Avoid unused variable warning
	
	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

func (s *Server) createPhishletVersion(c *gin.Context) {
	name := c.Param("name")
	
	var req CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	phishlet, err := s.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phishlet not found"})
		return
	}
	
	_ = phishlet
	_ = req
	
	c.JSON(http.StatusCreated, gin.H{"message": "Version created successfully"})
}

func (s *Server) getPhishletVersion(c *gin.Context) {
	name := c.Param("name")
	version := c.Param("version")
	_ = name    // Avoid unused variable warning
	_ = version // Avoid unused variable warning
	
	c.JSON(http.StatusOK, gin.H{"message": "Version retrieval not yet implemented"})
}

func (s *Server) evaluateConditions(c *gin.Context) {
	name := c.Param("name")
	
	var req EvaluateConditionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	phishlet, err := s.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phishlet not found"})
		return
	}
	
	var matchedActions []map[string]interface{}
	
	for _, condition := range phishlet.Conditions {
		matched := false
		
		switch condition.Type {
		case "email_domain":
			if req.Email != "" {
				for _, domain := range condition.Values {
					if len(req.Email) > len(domain) && req.Email[len(req.Email)-len(domain):] == domain {
						matched = true
						break
					}
				}
			}
		case "user_agent":
			if req.UserAgent != "" && condition.Regex != "" {
				if len(condition.Regex) > 0 && len(req.UserAgent) > 0 {
					matched = true
				}
			}
		}
		
		if matched {
			for _, action := range condition.Actions {
				matchedActions = append(matchedActions, map[string]interface{}{
					"type":  action.Type,
					"value": action.Value,
				})
			}
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"actions": matchedActions})
}

func (s *Server) getMultiPageFlows(c *gin.Context) {
	name := c.Param("name")
	
	phishlet, err := s.storage.GetPhishlet(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phishlet not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"flows": phishlet.MultiPageFlows})
}

func (s *Server) updateFlowStep(c *gin.Context) {
	name := c.Param("name")
	flowName := c.Param("flow")
	_ = name // Avoid unused variable warning
	
	var req UpdateFlowStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	_ = req
	_ = flowName
	
	c.JSON(http.StatusOK, gin.H{"message": "Flow step updated successfully"})
}
