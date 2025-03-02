package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/validators"
	"github.com/nicholas-fedor/eui64-calculator/ui"
)

// Calculator defines the interface for EUI-64 calculation logic used by handlers.
type Calculator interface {
	// CalculateEUI64 computes the EUI-64 interface ID and full IPv6 address from a MAC address and prefix.
	CalculateEUI64(mac, prefix string) (string, string, error)
}

// Handler manages HTTP request handling for the EUI-64 calculator application.
type Handler struct {
	calc Calculator // calc is the EUI-64 calculator implementation.
}

// NewHandler creates a new Handler with the specified EUI-64 calculator.
// It initializes the handler with the provided calculator for dependency injection.
func NewHandler(calc Calculator) *Handler {
	return &Handler{calc: calc}
}

// Home handles GET requests to the root path, rendering the home page.
// It serves the initial form for entering MAC and IPv6 prefix values, aborting with a 500 status on render failure.
func (h *Handler) Home(c *gin.Context) {
	if err := ui.Home().Render(c.Request.Context(), c.Writer); err != nil {
		slog.Error("Failed to render home page", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// Calculate handles POST requests to compute an EUI-64 address from form data.
// It validates the MAC address and IPv6 prefix from the request, computes the EUI-64 interface ID and full IPv6 address,
// and renders the result. Errors during validation or calculation are logged and displayed to the user.
func (h *Handler) Calculate(c *gin.Context) {
	mac := c.PostForm("mac")
	prefix := c.PostForm("ip-start")
	data := ui.ResultData{}

	if err := validators.ValidateMAC(mac); err != nil {
		data.Error = "Please enter a valid MAC address (e.g., 00-14-22-01-23-45)"

		slog.Warn("MAC validation failed", "mac", mac, "error", err)
		h.renderResult(c, data)

		return
	}

	if err := validators.ValidateIPv6Prefix(prefix); err != nil {
		data.Error = "Please enter a valid IPv6 prefix (e.g., 2001:db8::)"

		slog.Warn("Prefix validation failed", "prefix", prefix, "error", err)
		h.renderResult(c, data)

		return
	}

	interfaceID, fullIP, err := h.calc.CalculateEUI64(mac, prefix)
	data.InterfaceID = interfaceID
	data.FullIP = fullIP

	if err != nil {
		data.Error = "Failed to calculate EUI-64 address"

		slog.Error("EUI-64 calculation failed", "mac", mac, "prefix", prefix, "error", err)
	}

	h.renderResult(c, data)
}

// renderResult renders the calculation result to the HTTP response.
// It uses the provided ResultData to display either the computed EUI-64 address or an error message,
// aborting with a 500 status if rendering fails.
func (h *Handler) renderResult(c *gin.Context, data ui.ResultData) {
	if err := ui.Result(data).Render(c.Request.Context(), c.Writer); err != nil {
		slog.Error("Failed to render result", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
