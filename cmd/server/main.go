package main

import (
	"log"
	"os"
	"strings"

	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router, returning it for use.
func SetupRouter() *gin.Engine {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	r := gin.New()

	// Force log's color
	gin.ForceConsoleColor()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Configure trusted proxies
	trustedProxies := []string{}
	if proxies, exists := os.LookupEnv("TRUSTED_PROXIES"); exists && proxies != "" {
		trustedProxies = strings.Split(proxies, ",")
		for i, proxy := range trustedProxies {
			trustedProxies[i] = strings.TrimSpace(proxy)
		}
	}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Initialize handler
	handler := handlers.NewHandler()

	// Define routes
	r.GET("/", handler.Home)
	r.POST("/calculate", handler.Calculate)

	// Serve static files
	r.Static("/static", "./static")

	return r
}

func main() {
	router := SetupRouter()

	// Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
