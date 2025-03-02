package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter creates a Gin router for testing with the application’s configuration.
// It loads the server configuration and sets up the router, failing the test if either step encounters an error.
func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	config.StaticDir = filepath.Join("..", "..", "static")

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
			wantBody:   "body",
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
			router := setupRouter(t)

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

// TestTrustedProxies tests the router’s trusted proxy configuration.
// It ensures that client IP resolution works correctly with no proxies, a valid trusted proxy,
// and an invalid proxy format, verifying error handling in the latter case.
func TestTrustedProxies(t *testing.T) {
	tests := []struct {
		name           string
		trustedProxies string
		remoteAddr     string
		xForwardedFor  string
		wantClientIP   string
		wantStatus     int
		wantBody       string
		wantError      bool
	}{
		{
			name:           "No trusted proxies",
			trustedProxies: "",
			remoteAddr:     "192.168.1.1:12345",
			xForwardedFor:  "203.0.113.1",
			wantClientIP:   "192.168.1.1",
			wantStatus:     http.StatusOK,
			wantBody:       "EUI-64 Calculator",
		},
		{
			name:           "Trusted proxy",
			trustedProxies: "192.168.1.1",
			remoteAddr:     "192.168.1.1:12345",
			xForwardedFor:  "203.0.113.1",
			wantClientIP:   "203.0.113.1",
			wantStatus:     http.StatusOK,
			wantBody:       "EUI-64 Calculator",
		},
		{
			name:           "Invalid trusted proxy format",
			trustedProxies: "invalid-proxy,,192.168.1.1",
			remoteAddr:     "192.168.1.1:12345",
			xForwardedFor:  "203.0.113.1",
			wantClientIP:   "",
			wantStatus:     0,
			wantBody:       "",
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.trustedProxies != "" {
				t.Setenv(trustedProxiesEnv, tt.trustedProxies)
			} else {
				t.Setenv(trustedProxiesEnv, "")
			}

			if tt.wantError {
				_, err := SetupRouter(Config{Port: defaultPort, TrustedProxies: strings.Split(tt.trustedProxies, ",")})
				assert.Error(t, err, "Expected error for invalid proxy")

				return
			}

			router := setupRouter(t)
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remoteAddr

			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantStatus, resp.Code, "Status code")
			assert.Contains(t, resp.Body.String(), tt.wantBody, "Response body")

			var gotClientIP string

			router.GET("/test-ip", func(c *gin.Context) {
				gotClientIP = c.ClientIP()
			})

			req, _ = http.NewRequest(http.MethodGet, "/test-ip", nil)
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
