package api

import (
	"github.com/gin-gonic/gin"

	"go-wx/internal/handlers"
	"go-wx/internal/services"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())

	forecastHandler := handlers.NewForecastHandler(services.NewWeatherService())

	r.GET("/api/health", handlers.HealthCheck)
	r.GET("/api/forecast", forecastHandler.GetForecast)

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
