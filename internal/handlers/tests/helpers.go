package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers/mocks"
	"github.com/stretchr/testify/require"
)

// prepareCalcRequest sets up a POST request for Calculate handler tests.
func prepareCalcRequest(t *testing.T, formData url.Values) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	responseRecorder := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(responseRecorder)

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/calculate",
		strings.NewReader(formData.Encode()),
	)
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	err = req.ParseForm()
	require.NoError(t, err, "Failed to parse form")

	ginContext.Request = req

	return ginContext, responseRecorder
}

// setupInvalidHandler creates a handler for invalid Calculate tests.
func setupInvalidHandler(
	t *testing.T,
	validator *mocks.Validator,
	renderer *mocks.Renderer,
) (*handlers.Handler, *mocks.Renderer) {
	t.Helper()

	if renderer == nil {
		renderer = &mocks.Renderer{}
	}

	handler := handlers.NewHandler(&mocks.Calculator{}, validator, renderer)

	return handler, renderer
}
