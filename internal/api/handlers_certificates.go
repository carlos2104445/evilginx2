package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CertificateInfo struct {
	Domain     string `json:"domain"`
	Issuer     string `json:"issuer"`
	NotBefore  string `json:"not_before"`
	NotAfter   string `json:"not_after"`
	IsValid    bool   `json:"is_valid"`
	IsWildcard bool   `json:"is_wildcard"`
}

func (s *Server) listCertificates(c *gin.Context) {
	certificates := []CertificateInfo{
		{
			Domain:     "example.com",
			Issuer:     "Let's Encrypt",
			NotBefore:  "2024-01-01T00:00:00Z",
			NotAfter:   "2024-04-01T00:00:00Z",
			IsValid:    true,
			IsWildcard: false,
		},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"certificates": certificates,
		"count":        len(certificates),
	})
}

func (s *Server) generateCertificate(c *gin.Context) {
	var req struct {
		Domain string `json:"domain" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	cert := CertificateInfo{
		Domain:     req.Domain,
		Issuer:     "Let's Encrypt",
		NotBefore:  "2024-01-01T00:00:00Z",
		NotAfter:   "2024-04-01T00:00:00Z",
		IsValid:    true,
		IsWildcard: false,
	}
	
	c.JSON(http.StatusCreated, cert)
}

func (s *Server) deleteCertificate(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain is required"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "certificate deleted successfully"})
}
