package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers/mocks"
	"github.com/stretchr/testify/require"
)

func TestHomeHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		renderer   *mocks.MockRenderer // Use MockRenderer
		wantBody   string
		wantStatus int
	}{
		{
			name:       "Successful GET request",
			renderer:   &mocks.MockRenderer{},
			wantStatus: http.StatusOK,
			wantBody:   "EUI-64 Calculator",
		},
		{
			name:       "RenderHome failure",
			renderer:   &mocks.MockRenderer{HomeErr: errors.New("render home failed")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "", // No body expected on error
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			responseRecorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(responseRecorder)
			ginContext.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

			handler := handlers.NewHandler(
				&mocks.Calculator{InterfaceID: "", FullIP: "", Err: nil},
				&mocks.Validator{MacErr: nil, PrefixErr: nil},
				testCase.renderer,
			)
			handler.Home(mocks.NewRequestContext(ginContext)) // Use mocks.NewRequestContext
			require.Equal(t, testCase.wantStatus, responseRecorder.Code)

			if testCase.wantBody != "" {
				require.Contains(t, responseRecorder.Body.String(), testCase.wantBody)
			}

			require.True(t, testCase.renderer.CalledHome, "RenderHome not called")
		})
	}
}
