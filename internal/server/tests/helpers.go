package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// mockHandlerFunc simulates a gin.HandlerFunc for testing.
type mockHandlerFunc struct {
	called bool
}

func (m *mockHandlerFunc) ServeHTTP(c *gin.Context) {
	m.called = true

	c.String(http.StatusOK, "Mock Response")
}

// testRouterRoute performs a single route test for SetupRouter.
func testRouterRoute(
	t *testing.T,
	router *gin.Engine,
	httpMethod,
	path string,
	wStatus int,
	wBody string,
	handler *mockHandlerFunc,
	useHandlerCheck bool, // Added to control handler check
) {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), httpMethod, path, nil)
	require.NoError(t, err, "Failed to create request")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, wStatus, resp.Code)
	require.Contains(t, resp.Body.String(), wBody)

	if useHandlerCheck {
		if httpMethod == http.MethodGet && path == "/" {
			require.True(t, handler.called, "Home handler not called")
		} else if httpMethod == "POST" && path == "/calculate" {
			require.True(t, handler.called, "Calculate handler not called")
		}
	}
}
