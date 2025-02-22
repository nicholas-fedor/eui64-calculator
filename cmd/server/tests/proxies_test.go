package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/app"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/stretchr/testify/require"
)

// createProxyConfig initializes the server.Config for TestTrustedProxies without modifying env.
func createProxyConfig(t *testing.T, trustedProxies []string) server.Config {
	t.Helper()

	return server.Config{
		Port:           app.DefaultPort,
		TrustedProxies: trustedProxies,
		StaticDir:      "",
	}
}

// testProxySetup tests the proxy setup scenario.
func testProxySetup(t *testing.T, config server.Config) (*gin.Engine, error) {
	t.Helper()

	router, err := server.SetupRouter(config, func(c *gin.Context) {
		c.String(http.StatusOK, "EUI-64 Calculator")
	}, func(c *gin.Context) {
		c.String(http.StatusOK, "Calculate Response")
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set up proxy test router: %w", err)
	}

	return router, nil
}

// testProxyRequest tests a single proxy request and client IP.
func testProxyRequest(t *testing.T, router *gin.Engine, rAddr, xFwd, wantIP, wantBody string, wantStatus int) {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	require.NoError(t, err, "Failed to create request")

	req.RemoteAddr = rAddr
	if xFwd != "" {
		req.Header.Set("X-Forwarded-For", xFwd)
	}

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, wantStatus, resp.Code, "Status code")
	require.Contains(t, resp.Body.String(), wantBody, "Response body")

	var gotClientIP string

	router.GET("/test-ip", func(c *gin.Context) {
		gotClientIP = c.ClientIP()
	})

	req, err = http.NewRequestWithContext(t.Context(), http.MethodGet, "/test-ip", nil)
	require.NoError(t, err, "Failed to create test-ip request")

	req.RemoteAddr = rAddr
	if xFwd != "" {
		req.Header.Set("X-Forwarded-For", xFwd)
	}

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, wantIP, gotClientIP, "Client IP")
}

func TestTrustedProxies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		wantClientIP   string
		wantBody       string
		trustedProxies []string
		wantStatus     int
		wantError      bool
	}{
		{
			trustedProxies: nil,
			name:           "No trusted proxies",
			remoteAddr:     "192.168.1.1:12345",
			xForwardedFor:  "203.0.113.1",
			wantClientIP:   "192.168.1.1",
			wantBody:       "EUI-64 Calculator",
			wantError:      false,
			wantStatus:     http.StatusOK,
		},
		{
			trustedProxies: []string{"192.168.1.1"},
			name:           "Trusted proxy",
			remoteAddr:     "192.168.1.1:12345",
			xForwardedFor:  "203.0.113.1",
			wantClientIP:   "203.0.113.1",
			wantBody:       "EUI-64 Calculator",
			wantError:      false,
			wantStatus:     http.StatusOK,
		},
		{
			trustedProxies: []string{"invalid-proxy", "", "192.168.1.1"},
			name:           "Invalid trusted proxy format",
			remoteAddr:     "192.168.1.1:12345",
			xForwardedFor:  "203.0.113.1",
			wantClientIP:   "",
			wantBody:       "",
			wantError:      true,
			wantStatus:     0,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			runProxyTestCase(t, testCase)
		})
	}
}

// runProxyTestCase executes a single TestTrustedProxies case.
func runProxyTestCase(t *testing.T, testCase struct {
	name           string
	remoteAddr     string
	xForwardedFor  string
	wantClientIP   string
	wantBody       string
	trustedProxies []string
	wantStatus     int
	wantError      bool
}) {
	t.Helper()

	config := createProxyConfig(t, testCase.trustedProxies)

	if testCase.wantError {
		_, err := testProxySetup(t, config)
		require.Error(t, err, "Expected error for invalid proxy")

		return
	}

	router, err := testProxySetup(t, config)
	require.NoError(t, err)

	testProxyRequest(
		t,
		router,
		testCase.remoteAddr,
		testCase.xForwardedFor,
		testCase.wantClientIP,
		testCase.wantBody,
		testCase.wantStatus,
	)
}
