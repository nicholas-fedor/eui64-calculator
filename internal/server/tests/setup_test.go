package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/stretchr/testify/require"
)

func TestSetupRouter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		config     server.Config
		wantStatus int
		wantBody   string
		wantErr    bool
	}{
		{
			name: "Valid config with routes",
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: []string{"127.0.0.1"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "Mock Response",
			wantErr:    false,
		},
		{
			name: "Invalid trusted proxies",
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: []string{"invalid-proxy"},
			},
			wantStatus: 0,
			wantErr:    true,
		},
		{
			name: "Empty static dir",
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "",
				TrustedProxies: []string{"127.0.0.1"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "Mock Response",
			wantErr:    false,
		},
		{
			name: "No trusted proxies",
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: nil,
			},
			wantStatus: http.StatusOK,
			wantBody:   "Mock Response",
			wantErr:    false,
		},
		{
			name: "Default handlers",
			config: server.Config{
				Port:           ":8080",
				StaticDir:      "/tmp/static",
				TrustedProxies: []string{"127.0.0.1"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "EUI-64 Calculator",
			wantErr:    false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var homeHandler, calcHandler gin.HandlerFunc
			if testCase.name == "Default handlers" {
				homeHandler = nil
				calcHandler = nil
			} else {
				homeHandlerMock := &mockHandlerFunc{called: false}
				calcHandlerMock := &mockHandlerFunc{called: false}
				homeHandler = homeHandlerMock.ServeHTTP
				calcHandler = calcHandlerMock.ServeHTTP
			}

			router, err := server.SetupRouter(testCase.config, homeHandler, calcHandler)
			if testCase.wantErr {
				require.Error(t, err, "SetupRouter should fail with invalid config")
				require.Nil(t, router, "Router should be nil on error")

				return
			}

			require.NoError(t, err, "SetupRouter should succeed with valid config")
			require.NotNil(t, router, "Router should be initialized")

			// Test the home route
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			router.ServeHTTP(w, req)
			require.Equal(t, testCase.wantStatus, w.Code)

			if testCase.wantBody != "" {
				require.Contains(t, w.Body.String(), testCase.wantBody)
			}

			// Test calculate route
			w = httptest.NewRecorder()
			req, _ = http.NewRequest(http.MethodPost, "/calculate", nil)
			router.ServeHTTP(w, req)
			require.Equal(t, testCase.wantStatus, w.Code)

			if testCase.wantBody != "" {
				require.Contains(t, w.Body.String(), testCase.wantBody)
			}

			// Test static route (if StaticDir is set)
			if testCase.config.StaticDir != "" {
				w = httptest.NewRecorder()
				req, _ = http.NewRequest(http.MethodGet, "/static/test.css", nil)
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusNotFound, w.Code, "Static route should return 404 for missing file")

				// Test HEAD request for static route
				w = httptest.NewRecorder()
				req, _ = http.NewRequest(http.MethodHead, "/static/test.css", nil)
				router.ServeHTTP(w, req)
				require.Equal(t, http.StatusNotFound, w.Code, "HEAD request should return 404 for missing file")
			}
		})
	}
}
