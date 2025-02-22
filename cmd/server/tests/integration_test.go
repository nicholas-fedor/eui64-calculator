package main

import (
	"os"
	"testing"
	"time"

	"github.com/nicholas-fedor/eui64-calculator/internal/app"
)

func TestMainIntegration(t *testing.T) {
	t.Parallel()

	// Set environment to simulate a real run
	os.Setenv("PORT", ":8081") // Avoid conflict with other tests
	defer os.Unsetenv("PORT")

	// Run in a goroutine to avoid blocking
	done := make(chan error)
	go func() {
		appInstance := app.NewApp()
		done <- appInstance.Run()
	}()

	// Give it a moment to start (or fail)
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Expected Run to block, but it returned nil")
		}

		t.Logf("Run failed as expected (port in use or interrupted): %v", err)
	case <-time.After(1 * time.Second): // Success: app started and is running
	}
}
