package tests

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// mockHandlerFunc simulates a gin.HandlerFunc for testing.
type mockHandlerFunc struct {
	called bool
}

func (m *mockHandlerFunc) ServeHTTP(c *gin.Context) {
	m.called = true

	c.String(http.StatusOK, "Mock Response")
}
