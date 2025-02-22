package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/app"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // Set release mode for production
	gin.ForceConsoleColor()
	appInstance := app.NewApp()
	if err := appInstance.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
