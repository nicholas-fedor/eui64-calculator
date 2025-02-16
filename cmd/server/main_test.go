package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
	"github.com/stretchr/testify/assert"
)

// setupRouter creates a Gin router with the same configuration as SetupRouter, but adjusted for testing.
func setupRouter() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.New()

	// Force log's color (not relevant for testing, but included for completeness)
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
		panic(err) // In a real test, you might want to handle this error differently
	}

	// Initialize handler
	handler := handlers.NewHandler()

	// Define routes
	r.GET("/", handler.Home)
	r.POST("/calculate", handler.Calculate)

	// Serve static files - use an absolute path to the static directory
	_, filename, _, _ := runtime.Caller(0)                        // Get the current file's path
	projectRoot := filepath.Join(filepath.Dir(filename), "../..") // Go up two levels to the project root
	staticPath := filepath.Join(projectRoot, "static")            // Absolute path to static directory
	r.Static("/static", staticPath)

	return r
}

// TestRouterSetup tests the router configuration and route behavior.
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
			wantBody:   "body", // Check for some content from styles.css
		},
		{
			name:       "GET /nonexistent - Not found",
			method:     "GET",
			path:       "/nonexistent",
			wantStatus: http.StatusNotFound,
			wantBody:   "404 page not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			var req *http.Request
			if tt.method == "POST" {
				req, _ = http.NewRequest(tt.method, tt.path, strings.NewReader(tt.formData.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			} else {
				req, _ = http.NewRequest(tt.method, tt.path, nil)
			}
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Status code")
			assert.Contains(t, resp.Body.String(), tt.wantBody, "Response body")
		})
	}
}

// TestTrustedProxies tests the trusted proxies configuration.
func TestTrustedProxies(t *testing.T) {
	tests := []struct {
		name           string
		trustedProxies string
		remoteAddr     string
		xForwardedFor  string
		wantClientIP   string
		wantStatus     int
		wantBody       string
	}{
		{
			name:           "No trusted proxies",
			trustedProxies: "",
			remoteAddr:     "192.168.1.1:12345", // Include a port number
			xForwardedFor:  "203.0.113.1",
			wantClientIP:   "192.168.1.1", // Should use remote address, not X-Forwarded-For
			wantStatus:     http.StatusOK,
			wantBody:       "EUI-64 Calculator",
		},
		{
			name:           "Trusted proxy",
			trustedProxies: "192.168.1.1",
			remoteAddr:     "192.168.1.1:12345", // Include a port number
			xForwardedFor:  "203.0.113.1",
			wantClientIP:   "203.0.113.1", // Should use X-Forwarded-For
			wantStatus:     http.StatusOK,
			wantBody:       "EUI-64 Calculator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set TRUSTED_PROXIES environment variable for this test
			if tt.trustedProxies != "" {
				t.Setenv("TRUSTED_PROXIES", tt.trustedProxies)
			}

			router := setupRouter()

			// Test the home route to ensure the router is working
			req, _ := http.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Status code")
			assert.Contains(t, resp.Body.String(), tt.wantBody, "Response body")

			// Test client IP by adding a temporary route
			var gotClientIP string
			router.GET("/test-ip", func(c *gin.Context) {
				gotClientIP = c.ClientIP()
			})

			req, _ = http.NewRequest("GET", "/test-ip", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			resp = httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantClientIP, gotClientIP, "Client IP")
		})
	}
}
