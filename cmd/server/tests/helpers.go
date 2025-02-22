package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/app"
	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/nicholas-fedor/eui64-calculator/internal/validators"
	"github.com/stretchr/testify/require"
)

// setupRouter creates a Gin router for routes testing with a simplified configuration.
func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()

	config := server.Config{
		Port:           app.DefaultPort,
		StaticDir:      "../../../static", // Adjusted to reach project root static/ from cmd/server/tests/
		TrustedProxies: []string{},
	}
	calculator := &eui64.DefaultCalculator{}
	validator := &validators.CombinedValidator{}
	handler := handlers.NewHandler(calculator, validator, &server.UIRenderer{})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	err := router.SetTrustedProxies(config.TrustedProxies)
	require.NoError(t, err, "Failed to set trusted proxies")

	router.GET("/", handler.HomeAdapter())
	router.POST("/calculate", handler.CalculateAdapter())
	router.Static("/static", config.StaticDir)

	return router
}

// prepareRouteRequest prepares an HTTP request for testing.
func prepareRouteRequest(t *testing.T, method, path string, formData url.Values) *http.Request {
	t.Helper()

	var request *http.Request

	var err error

	if method == "POST" {
		request, err = http.NewRequestWithContext(t.Context(), method, path, strings.NewReader(formData.Encode()))
		require.NoError(t, err, "Failed to create POST request")
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		request, err = http.NewRequestWithContext(t.Context(), method, path, nil)
		require.NoError(t, err, "Failed to create GET request")
	}

	return request
}

// testRoute performs a single route test.
func testRoute(t *testing.T, router *gin.Engine, request *http.Request, wantStatus int, wantBody string) {
	t.Helper()

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	require.Equal(t, wantStatus, responseRecorder.Code)
	require.Contains(t, responseRecorder.Body.String(), wantBody)
}
