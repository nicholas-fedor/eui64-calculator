package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/EUI64-Calculator/internal/eui64"
	"github.com/nicholas-fedor/EUI64-Calculator/ui"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Home(c *gin.Context) {
	if err := ui.Home().Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) Calculate(c *gin.Context) {
	mac := c.PostForm("mac")
	prefix := c.PostForm("ip-start")

	interfaceID, fullIP, err := eui64.CalculateEUI64(mac, prefix)
	data := ui.ResultData{
		InterfaceID: interfaceID,
		FullIP:      fullIP,
		Error:       "",
	}
	if err != nil {
		data.Error = err.Error()
	}

	if err := ui.Result(data).Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
