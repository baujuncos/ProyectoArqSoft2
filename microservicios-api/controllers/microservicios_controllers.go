package controllers

import (
	"github.com/gin-gonic/gin"
	"microservicios-api/services"
	"net/http"
)

func GetServices(c *gin.Context) {
	services, err := services.GetServices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}
