package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-wx/internal/handlers"
	"go-wx/internal/services"
)

const (
	HealthRoute  = "/api/health"
	ForecastRoute = "/api/forecast"

	CORSAllowOrigin  = "*"
	CORSAllowMethods = "GET, POST, PUT, DELETE, OPTIONS"
	CORSAllowHeaders = "Content-Type, Authorization"

	HeaderAllowOrigin  = "Access-Control-Allow-Origin"
	HeaderAllowMethods = "Access-Control-Allow-Methods"
	HeaderAllowHeaders = "Access-Control-Allow-Headers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())

	forecastHandler := handlers.NewForecastHandler(services.NewWeatherService())

	r.GET(HealthRoute, handlers.HealthCheck)
	r.GET(ForecastRoute, forecastHandler.GetForecast)

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(HeaderAllowOrigin, CORSAllowOrigin)
		c.Header(HeaderAllowMethods, CORSAllowMethods)
		c.Header(HeaderAllowHeaders, CORSAllowHeaders)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
