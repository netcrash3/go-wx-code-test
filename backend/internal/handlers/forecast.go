package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-wx/internal/models"
)

const (
	QueryParamLat = "lat"
	QueryParamLon = "lon"

	ErrMissingParams = "lat and lon query parameters are required"
	ErrInvalidLat    = "lat must be a valid number"
	ErrInvalidLon    = "lon must be a valid number"

	ResponseKeyError = "error"
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
	latStr := c.Query(QueryParamLat)
	lonStr := c.Query(QueryParamLon)

	if latStr == "" || lonStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{ResponseKeyError: ErrMissingParams})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ResponseKeyError: ErrInvalidLat})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ResponseKeyError: ErrInvalidLon})
		return
	}

	result, err := h.service.GetTodayForecast(lat, lon)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{ResponseKeyError: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
