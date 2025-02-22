package tests

import (
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
		wantBody   string
		wantStatus int
	}{
		{
			name:       "Successful GET request",
			wantStatus: http.StatusOK,
			wantBody:   "EUI-64 Calculator",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			responseRecorder := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(responseRecorder)
			ginContext.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

			renderer := &mocks.Renderer{
				HomeErr:      nil,
				ResultErr:    nil,
				CalledHome:   false,
				CalledResult: false,
			}
			handler := handlers.NewHandler(
				&mocks.Calculator{InterfaceID: "", FullIP: "", Err: nil},
				&mocks.Validator{MacErr: nil, PrefixErr: nil},
				renderer,
			)
			ctx := mocks.NewRequestContext(ginContext) // Use constructor
			handler.Home(ctx)
			require.Equal(t, testCase.wantStatus, responseRecorder.Code)
			require.Contains(t, responseRecorder.Body.String(), testCase.wantBody)
			require.True(t, renderer.CalledHome, "RenderHome not called")
		})
	}
}
