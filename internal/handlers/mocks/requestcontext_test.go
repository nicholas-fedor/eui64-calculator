package mocks

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("FormValue", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Add form data
		c.Request, _ = http.NewRequest(http.MethodPost, "/", nil)
		c.Request.PostFormValue("testkey") // Sets form value

		rc := NewRequestContext(c)
		value := rc.FormValue("testkey")

		if value != "" { // Expect empty since we didn't actually set a value
			t.Errorf("FormValue() = %v, want empty string", value)
		}
	})

	t.Run("GetContext", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		rc := NewRequestContext(c)
		gotContext := rc.GetContext()

		if gotContext != c {
			t.Errorf("GetContext() = %v, want %v", gotContext, c)
		}
	})
}
