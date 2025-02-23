package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers/mocks"
	"github.com/stretchr/testify/require"
)

func TestRenderResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                          string
		renderer                      *mocks.MockRenderer
		status                        int
		interfaceID, fullIP, errorMsg string
		wantStatus                    int
		wantBody                      string
	}{
		{
			name:        "Success render",
			renderer:    &mocks.MockRenderer{},
			status:      http.StatusOK,
			interfaceID: "0214:22ff:fe01:2345",
			fullIP:      "2001:db8::214:22ff:fe01:2345",
			errorMsg:    "",
			wantStatus:  http.StatusOK,
			wantBody:    "0214:22ff:fe01:2345",
		},
		{
			name:        "Render failure",
			renderer:    &mocks.MockRenderer{ResultErr: errors.New("render failed")},
			status:      http.StatusOK,
			interfaceID: "",
			fullIP:      "",
			errorMsg:    "",
			wantStatus:  http.StatusInternalServerError,
			wantBody:    "render failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Status(tt.status)
			h := NewHandler(&mocks.Calculator{}, &mocks.Validator{}, tt.renderer)
			h.renderResult(c, tt.status, tt.interfaceID, tt.fullIP, tt.errorMsg)
			require.Equal(t, tt.wantStatus, w.Code)
			require.Contains(t, w.Body.String(), tt.wantBody)
		})
	}
}
