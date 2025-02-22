package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/nicholas-fedor/eui64-calculator/ui"
	"github.com/stretchr/testify/require"
)

func TestUIRenderer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		data       ui.ResultData
		name       string
		wantBody   string
		wantStatus int
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
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			renderer := &server.UIRenderer{}
			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			ctx.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

			isHome := testCase.name == "RenderHome"
			if isHome {
				err := renderer.RenderHome(ctx)
				require.NoError(t, err)
			} else {
				err := renderer.RenderResult(ctx, testCase.data.InterfaceID, testCase.data.FullIP, testCase.data.Error)
				require.NoError(t, err)
			}

			require.Equal(t, testCase.wantStatus, resp.Code)
			require.Contains(t, resp.Body.String(), testCase.wantBody)
		})
	}
}
