package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MockRenderer struct {
	HomeErr      error
	ResultErr    error
	CalledHome   bool
	CalledResult bool
}

func (m *MockRenderer) RenderHome(ctx *gin.Context) error {
	m.CalledHome = true
	if m.HomeErr != nil {
		return m.HomeErr
	}

	ctx.String(http.StatusOK, "EUI-64 Calculator")

	return nil
}

func (m *MockRenderer) RenderResult(ctx *gin.Context, interfaceID, fullIP, errorMsg string) error {
	m.CalledResult = true
	if m.ResultErr != nil {
		return m.ResultErr
	}
	// Use the current status from the context instead of hardcoding 200
	status := ctx.Writer.Status()
	if status == 0 { // Default to 200 if not set
		status = http.StatusOK
	}

	if errorMsg != "" {
		ctx.String(status, errorMsg)
	} else {
		ctx.String(status, "%s\n%s", interfaceID, fullIP)
	}

	return nil
}
