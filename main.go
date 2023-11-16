package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/skanderphilipp/sisyApi/internal"
)

func main() {
	// Initialize the application with Wire
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	router := gin.Default()
	// Setup routes
	setupRoutes(router, app)

	// Start the application
	if err := router.Run(":8080"); err != nil {
		// handle error, perhaps log it and/or gracefully shut down
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRoutes(router *gin.Engine, app *internal.App) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
