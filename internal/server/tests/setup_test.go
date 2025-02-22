package tests

import (
	"net/http"
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/stretchr/testify/require"
)

func TestSetupRouter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		config     server.Config
		name       string
		wantBody   string
		wantStatus int
		wantError  bool
	}{
		{
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: []string{"192.168.1.1"},
			},
			name:       "Valid config with routes",
			wantBody:   "Mock Response",
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: []string{"invalid-proxy"},
			},
			name:       "Invalid trusted proxies",
			wantBody:   "",
			wantStatus: 0,
			wantError:  true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			homeHandler := &mockHandlerFunc{called: false}
			calcHandler := &mockHandlerFunc{called: false}

			router, err := server.SetupRouter(testCase.config, homeHandler.ServeHTTP, calcHandler.ServeHTTP)
			if testCase.wantError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, router)

			testRouterRoute(t, router, http.MethodGet, "/", testCase.wantStatus, testCase.wantBody, homeHandler)
			testRouterRoute(t, router, "POST", "/calculate", testCase.wantStatus, testCase.wantBody, calcHandler)
			testRouterRoute(t, router, http.MethodGet, "/static/test.css", http.StatusNotFound, "", nil)
		})
	}
}
