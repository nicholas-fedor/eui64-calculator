package mocks

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMockRenderer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		homeErr   error
		resultErr error
	}{
		{"No errors", nil, nil},
		{"Home error", errors.New("home error"), nil},
		{"Result error", nil, errors.New("result error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockRenderer{
				HomeErr:   tt.homeErr,
				ResultErr: tt.resultErr,
			}

			// Test RenderHome
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			err := m.RenderHome(c)

			if !errors.Is(err, tt.homeErr) {
				t.Errorf("RenderHome() error = %v, want %v", err, tt.homeErr)
			}

			if !m.CalledHome {
				t.Error("RenderHome() didn't set CalledHome")
			}

			// Test RenderResult
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			err = m.RenderResult(c, "interface", "full_ip", "")

			if !errors.Is(err, tt.resultErr) {
				t.Errorf("RenderResult() error = %v, want %v", err, tt.resultErr)
			}

			if !m.CalledResult {
				t.Error("RenderResult() didn't set CalledResult")
			}

			if tt.resultErr == nil && w.Code != http.StatusOK {
				t.Errorf("RenderResult() status = %v, want %v", w.Code, http.StatusOK)
			}
		})
	}
}
