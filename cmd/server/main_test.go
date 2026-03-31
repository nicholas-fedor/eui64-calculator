package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRouter creates a Fiber app for testing with the application's configuration.
// It loads the server configuration and sets up the app, failing the test if either step encounters an error.
func setupRouter(t *testing.T) *fiber.App {
	t.Helper()

	config := LoadConfig()

	app, err := SetupRouter(config)
	if err != nil {
		t.Fatalf("Failed to setup router: %v", err)
	}

	return app
}

// TestRouterSetup tests the app's handling of various HTTP requests.
// It verifies that the app correctly serves the home page, handles valid and invalid EUI-64 calculation requests,
// serves static files, and returns a 404 for unknown paths.
func TestRouterSetup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		formData   url.Values
		wantStatus int
		wantBody   string
	}{
		{
			name:       "GET / - Home page",
			method:     "GET",
			path:       "/",
			wantStatus: http.StatusOK,
			wantBody:   "EUI-64 Calculator",
		},
		{
			name:   "POST /calculate - Valid MAC and full prefix",
			method: "POST",
			path:   "/calculate",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "0214:22ff:fe01:2345",
		},
		{
			name:   "POST /calculate - Invalid MAC",
			method: "POST",
			path:   "/calculate",
			formData: url.Values{
				"mac":      {"invalid-mac"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "error-message",
		},
		{
			name:       "GET /static/styles.css - Static file",
			method:     "GET",
			path:       "/static/styles.css",
			wantStatus: http.StatusOK,
			wantBody:   "body {", // Partial match for CSS content; adjust if needed
		},
		{
			name:       "GET /unknown - 404 Not Found",
			method:     "GET",
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
		},
	}

	app := setupRouter(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var req *http.Request
			if tt.method == "POST" {
				req, _ = http.NewRequestWithContext(
					context.Background(),
					tt.method,
					tt.path,
					strings.NewReader(tt.formData.Encode()),
				)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else {
				req, _ = http.NewRequestWithContext(context.Background(), tt.method, tt.path, http.NoBody)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.wantStatus, resp.StatusCode, "Status code")

			if tt.wantBody != "" {
				assert.Contains(t, string(body), tt.wantBody, "Body content")
			}
		})
	}
}

// TestLoadConfig to cover Lines 37 and 56.
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		portEnv        string
		trustedProxies string
		wantPort       string
		wantProxies    []string
	}{
		{
			name:           "Default config",
			portEnv:        "",
			trustedProxies: "",
			wantPort:       ":" + defaultPort,
			wantProxies:    nil,
		},
		{
			name:           "Custom port (Line 37)",
			portEnv:        "9090",
			trustedProxies: "",
			wantPort:       ":9090",
			wantProxies:    nil,
		},
		{
			name:           "Trusted proxies with empty entry",
			portEnv:        "",
			trustedProxies: "192.168.1.1, ,192.168.1.2",
			wantPort:       ":" + defaultPort,
			wantProxies:    []string{"192.168.1.1", "192.168.1.2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear and set environment variables
			t.Setenv("PORT", tt.portEnv)
			t.Setenv(trustedProxiesEnv, tt.trustedProxies)

			config := LoadConfig()

			assert.Equal(t, tt.wantPort, config.Port, "Port")
			assert.Equal(t, tt.wantProxies, config.TrustedProxies, "TrustedProxies")
		})
	}
}
