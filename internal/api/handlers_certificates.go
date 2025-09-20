package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CertificateRequest struct {
	Domain string `json:"domain"`
	Cert   string `json:"cert"`
	Key    string `json:"key"`
}

func (h *Handlers) listCertificates(c *gin.Context) {
	certificates, err := h.storage.List(c.Request.Context(), "cert:")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"certificates": certificates})
}

func (h *Handlers) createCertificate(c *gin.Context) {
	var req CertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	key := "cert:" + req.Domain
	value := req.Cert + "|" + req.Key
	
	err := h.storage.Set(c.Request.Context(), key, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Certificate created successfully"})
}

func (h *Handlers) getCertificate(c *gin.Context) {
	domain := c.Param("domain")
	key := "cert:" + domain
	
	value, err := h.storage.Get(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"domain": domain, "data": value})
}

func (h *Handlers) deleteCertificate(c *gin.Context) {
	domain := c.Param("domain")
	key := "cert:" + domain
	
	err := h.storage.Delete(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted successfully"})
}
