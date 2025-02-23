package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/app"
)

// Package-level variables to allow mocking in tests.
var (
	newAppFunc = app.NewApp
	logFatalf  = log.Fatalf
	runFunc    = func() error { // Define run as a variable
		gin.SetMode(gin.ReleaseMode)
		gin.ForceConsoleColor()

		appInstance := newAppFunc()

		return appInstance.Run()
	}
)

func main() {
	if err := runFunc(); err != nil {
		logFatalf("Failed to start server: %v", err)
	}
}
