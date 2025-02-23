package server

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/stretchr/testify/require"
)

func TestSetupRouter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		config          server.Config
		name            string
		wantBody        string
		wantStatus      int
		wantError       bool
		useHandlerCheck bool // Added to control handler check
	}{
		{
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: []string{"192.168.1.1"},
			},
			name:            "Valid config with routes",
			wantBody:        "Mock Response",
			wantStatus:      http.StatusOK,
			wantError:       false,
			useHandlerCheck: true,
		},
		{
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: []string{"invalid-proxy"},
			},
			name:            "Invalid trusted proxies",
			wantBody:        "",
			wantStatus:      0,
			wantError:       true,
			useHandlerCheck: false,
		},
		{
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: []string{},
			},
			name:            "Default handlers",
			wantBody:        "EUI-64 Calculator",
			wantStatus:      http.StatusOK,
			wantError:       false,
			useHandlerCheck: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			homeHandler := &mockHandlerFunc{called: false}
			calcHandler := &mockHandlerFunc{called: false}

			var router *gin.Engine

			var err error
			if testCase.name == "Default handlers" {
				router, err = server.SetupRouter(testCase.config, nil, nil) // Test default handlers
			} else {
				router, err = server.SetupRouter(testCase.config, homeHandler.ServeHTTP, calcHandler.ServeHTTP)
			}

			if testCase.wantError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, router)

			testRouterRoute(
				t, router, http.MethodGet, "/", testCase.wantStatus, testCase.wantBody, homeHandler, testCase.useHandlerCheck,
			)
			testRouterRoute(
				t, router, "POST", "/calculate", testCase.wantStatus, testCase.wantBody, calcHandler, testCase.useHandlerCheck,
			)
			testRouterRoute(
				t, router, http.MethodGet, "/static/test.css", http.StatusNotFound, "", nil, false,
			)
		})
	}
}
