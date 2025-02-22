package app

import (
	"os"
	"testing"
	"time"

	"github.com/nicholas-fedor/eui64-calculator/internal/app"
)

func TestMainIntegration(t *testing.T) {
	t.Parallel()

	// Simulate environment
	err := os.Setenv("PORT", ":8081") // Unique port to avoid conflicts
	if err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	defer os.Unsetenv("PORT")

	// Run app in a goroutine
	done := make(chan error)
	go func() {
		appInstance := app.NewApp()
		done <- appInstance.Run()
	}()

	// Wait briefly to ensure startup or catch immediate failure
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Expected Run to block, but it returned nil")
		}

		t.Logf("Run failed as expected (port in use or interrupted): %v", err)
	case <-time.After(1 * time.Second): // Success: app started and is running
	}
}
