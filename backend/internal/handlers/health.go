package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-wx/internal/models"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{Status: "ok"})
}
