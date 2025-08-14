package main

import (
	// Added for context support.
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRouter creates a Gin router for testing with the application’s configuration.
// It loads the server configuration and sets up the router, failing the test if either step encounters an error.
func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()

	config := LoadConfig()

	r, err := SetupRouter(config)
	if err != nil {
		t.Fatalf("Failed to setup router: %v", err)
	}

	return r
}

// TestRouterSetup tests the router’s handling of various HTTP requests.
// It verifies that the router correctly serves the home page, handles valid and invalid EUI-64 calculation requests,
// serves static files, and returns a 404 for unknown paths.
func TestRouterSetup(t *testing.T) {
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

	router := setupRouter(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == "POST" {
				req = httptest.NewRequest(
					tt.method,
					tt.path,
					strings.NewReader(tt.formData.Encode()),
				)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Status code")

			if tt.wantBody != "" {
				assert.Contains(t, resp.Body.String(), tt.wantBody, "Body content")
			}
		})
	}
}

// TestTrustedProxies tests the router’s handling of trusted proxies.
// It verifies that the client IP is correctly determined based on the X-Forwarded-For header
// when trusted proxies are configured, covering various proxy scenarios.
func TestTrustedProxies(t *testing.T) {
	tests := []struct {
		name           string
		trustedProxies []string
		remoteAddr     string
		xForwardedFor  string
		wantClientIP   string
	}{
		{
			name:           "No proxies",
			trustedProxies: nil,
			remoteAddr:     "192.168.1.1:12345",
			xForwardedFor:  "",
			wantClientIP:   "192.168.1.1",
		},
		{
			name:           "Single trusted proxy",
			trustedProxies: []string{"192.168.1.2"},
			remoteAddr:     "192.168.1.2:54321",
			xForwardedFor:  "192.168.1.1",
			wantClientIP:   "192.168.1.1",
		},
		{
			name:           "Multiple trusted proxies",
			trustedProxies: []string{"192.168.1.3", "192.168.1.2"},
			remoteAddr:     "192.168.1.3:54321",
			xForwardedFor:  "192.168.1.1, 192.168.1.2",
			wantClientIP:   "192.168.1.1",
		},
		{
			name:           "Untrusted proxy",
			trustedProxies: []string{"192.168.1.2"},
			remoteAddr:     "192.168.1.3:54321",
			xForwardedFor:  "192.168.1.1",
			wantClientIP:   "192.168.1.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Port:           ":8080",
				TrustedProxies: tt.trustedProxies,
			}

			router, err := SetupRouter(config)
			require.NoError(t, err)

			// Set up a test route to capture the client IP.
			router.GET("/test-ip", func(c *gin.Context) {
				gotClientIP := c.ClientIP()
				c.String(http.StatusOK, gotClientIP)
			})

			req := httptest.NewRequest(http.MethodGet, "/test-ip", nil)
			req.RemoteAddr = tt.remoteAddr

			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantClientIP, strings.TrimSpace(resp.Body.String()), "Client IP")
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
