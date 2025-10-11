package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgretzky/evilginx2/pkg/models"
)

func (s *Server) listCertificates(c *gin.Context) {
	items, err := s.storage.ListCertificates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"certificates": items,
		"count":        len(items),
	})
}

func (s *Server) generateCertificate(c *gin.Context) {
	var req struct {
		Domain     string `json:"domain" binding:"required"`
		Issuer     string `json:"issuer"`
		NotBefore  string `json:"not_before"`
		NotAfter   string `json:"not_after"`
		IsWildcard bool   `json:"is_wildcard"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var nb, na time.Time
	var err error
	if req.NotBefore != "" {
		nb, _ = time.Parse(time.RFC3339, req.NotBefore)
	}
	if req.NotAfter != "" {
		na, _ = time.Parse(time.RFC3339, req.NotAfter)
	}

	cert := &models.Certificate{
		Domain:     req.Domain,
		Issuer:     req.Issuer,
		NotBefore:  nb,
		NotAfter:   na,
		IsValid:    na.IsZero() || na.After(time.Now()),
		IsWildcard: req.IsWildcard,
	}
	if err = s.storage.CreateCertificate(c.Request.Context(), cert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cert)
}

func (s *Server) deleteCertificate(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain is required"})
		return
	}
	if err := s.storage.DeleteCertificate(c.Request.Context(), domain); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "certificate deleted successfully"})
}
