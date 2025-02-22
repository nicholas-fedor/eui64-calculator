package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func run() error {
	gin.SetMode(gin.ReleaseMode) // L11-12
	gin.ForceConsoleColor()

	appInstance := app.NewApp() // L14

	return appInstance.Run() // L16-17
}
