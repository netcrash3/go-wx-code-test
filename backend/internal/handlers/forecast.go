package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-wx/internal/models"
)

type Forecaster interface {
	GetTodayForecast(lat, lon float64) (*models.ForecastResponse, error)
}

type ForecastHandler struct {
	service Forecaster
}

func NewForecastHandler(svc Forecaster) *ForecastHandler {
	return &ForecastHandler{service: svc}
}

func (h *ForecastHandler) GetForecast(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon query parameters are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat must be a valid number"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lon must be a valid number"})
		return
	}

	result, err := h.service.GetTodayForecast(lat, lon)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
