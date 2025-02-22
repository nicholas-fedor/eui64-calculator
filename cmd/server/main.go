package main

import (
	"log"

	"github.com/nicholas-fedor/eui64-calculator/internal/app"
)

func main() {
	appInstance := app.NewApp()
	if err := appInstance.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
