package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"go-wx/internal/api"
)

const (
	EnvPort     = "PORT"
	DefaultPort = "8080"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv(EnvPort)
	if port == "" {
		port = DefaultPort
	}

	router := api.NewRouter()

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
