package handlers

import (
	"log"
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
		if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
			// Log the error since there's nothing more we can do here
			log.Printf("Failed to abort with error: %v", err)
		}
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
		if err := c.AbortWithError(http.StatusInternalServerError, err); err != nil {
			// Log the error for debugging purposes
			log.Printf("Failed to abort with error: %v", err)
		}
		return
	}
}
