package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Renderer struct {
	HomeErr      error
	ResultErr    error
	CalledHome   bool
	CalledResult bool
}

func (m *Renderer) RenderHome(ctx *gin.Context) error {
	m.CalledHome = true

	ctx.String(http.StatusOK, "EUI-64 Calculator")

	return m.HomeErr
}

func (m *Renderer) RenderResult(ctx *gin.Context, interfaceID, fullIP, errorMsg string) error {
	m.CalledResult = true

	if errorMsg != "" {
		ctx.String(http.StatusOK, errorMsg)
	} else {
		ctx.String(http.StatusOK, "%s\n%s", interfaceID, fullIP)
	}

	return m.ResultErr
}
