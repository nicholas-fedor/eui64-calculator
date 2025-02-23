package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/nicholas-fedor/eui64-calculator/ui"
	"github.com/stretchr/testify/require"
)

// mockableRenderer wraps UIRenderer to allow error simulation.
type mockableRenderer struct {
	*server.UIRenderer
	forceError bool
}

func (r *mockableRenderer) RenderHome(ctx *gin.Context) error {
	if r.forceError {
		return errors.New("simulated render failure")
	}

	return r.UIRenderer.RenderHome(ctx)
}

func (r *mockableRenderer) RenderResult(ctx *gin.Context, interfaceID, fullIP, errorMsg string) error {
	if r.forceError {
		return errors.New("simulated render failure")
	}

	return r.UIRenderer.RenderResult(ctx, interfaceID, fullIP, errorMsg)
}

func TestUIRenderer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		data       ui.ResultData
		name       string
		wantBody   string
		wantStatus int
		isHome     bool
		forceError bool // Flag to simulate rendering error
	}{
		{
			data: ui.ResultData{
				InterfaceID: "",
				FullIP:      "",
				Error:       "",
			},
			name:       "RenderHome",
			wantBody:   "EUI-64 Calculator",
			wantStatus: http.StatusOK,
			isHome:     true,
			forceError: false,
		},
		{
			data: ui.ResultData{
				InterfaceID: "",
				FullIP:      "",
				Error:       "Invalid input",
			},
			name:       "RenderResult with error",
			wantBody:   "Invalid input",
			wantStatus: http.StatusOK,
			isHome:     false,
			forceError: false,
		},
		{
			data: ui.ResultData{
				InterfaceID: "0214:22ff:fe01:2345",
				FullIP:      "2001:db8::214:22ff:fe01:2345",
				Error:       "",
			},
			name:       "RenderResult with valid data",
			wantBody:   "0214:22ff:fe01:2345",
			wantStatus: http.StatusOK,
			isHome:     false,
			forceError: false,
		},
		{
			data:       ui.ResultData{},
			name:       "RenderHome_error",
			wantBody:   "",
			wantStatus: http.StatusInternalServerError,
			isHome:     true,
			forceError: true,
		},
		{
			data:       ui.ResultData{},
			name:       "RenderResult_error",
			wantBody:   "",
			wantStatus: http.StatusInternalServerError,
			isHome:     false,
			forceError: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Use mockableRenderer to control error simulation
			renderer := &mockableRenderer{
				UIRenderer: &server.UIRenderer{},
				forceError: testCase.forceError,
			}
			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			ctx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

			var err error
			if testCase.isHome {
				err = renderer.RenderHome(ctx)
			} else {
				err = renderer.RenderResult(ctx, testCase.data.InterfaceID, testCase.data.FullIP, testCase.data.Error)
			}

			// Simulate handler behavior: set status based on error
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			if testCase.forceError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, testCase.wantStatus, resp.Code)

			if testCase.wantBody != "" {
				require.Contains(t, resp.Body.String(), testCase.wantBody)
			}
		})
	}
}
