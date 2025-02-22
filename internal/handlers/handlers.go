package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Validator defines the interface for validating input data.
type Validator interface {
	ValidateMAC(mac string) error
	ValidateIPv6Prefix(prefix string) error
}

// RequestContext abstracts HTTP request handling using Gin’s context.
type RequestContext interface {
	FormValue(key string) string
	GetContext() *gin.Context
}

// Calculator defines the interface for EUI-64 calculation logic.
type Calculator interface {
	CalculateEUI64(mac, prefix string) (string, string, error)
}

// Renderer defines the interface for rendering UI components.
type Renderer interface {
	RenderHome(ctx *gin.Context) error
	RenderResult(ctx *gin.Context, interfaceID, fullIP, errorMsg string) error
}

// Handler manages HTTP request handling for the EUI-64 calculator.
type Handler struct {
	calc      Calculator
	validator Validator
	renderer  Renderer
}

// NewHandler creates a new Handler with injected dependencies.
func NewHandler(calc Calculator, validator Validator, renderer Renderer) *Handler {
	return &Handler{
		calc:      calc,
		validator: validator,
		renderer:  renderer,
	}
}

// Home handles GET requests to the root path, rendering the home page.
func (h *Handler) Home(ctx RequestContext) {
	c := ctx.GetContext()
	if err := h.renderer.RenderHome(c); err != nil {
		slog.Error("Failed to render home page", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Status(http.StatusOK)
}

// HomeAdapter adapts the Home method to a gin.HandlerFunc.
func (h *Handler) HomeAdapter() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.Home(&ginRequestContext{c: c})
	}
}

// Calculate handles POST requests to compute an EUI-64 address from form data.
func (h *Handler) Calculate(c *gin.Context) {
	mac := c.PostForm("mac")
	prefix := c.PostForm("ip-start")
	interfaceID, fullIP, errorMsg := "", "", ""

	if err := h.validator.ValidateMAC(mac); err != nil {
		errorMsg = "Please enter a valid MAC address (e.g., 00-14-22-01-23-45)"

		slog.Warn("MAC validation failed", "mac", mac, "error", err)
		h.renderResult(c, http.StatusOK, interfaceID, fullIP, errorMsg)

		return
	}

	if err := h.validator.ValidateIPv6Prefix(prefix); err != nil {
		errorMsg = "Please enter a valid IPv6 prefix (e.g., 2001:db8::)"

		slog.Warn("Prefix validation failed", "prefix", prefix, "error", err)
		h.renderResult(c, http.StatusOK, interfaceID, fullIP, errorMsg)

		return
	}

	interfaceID, fullIP, err := h.calc.CalculateEUI64(mac, prefix)
	if err != nil {
		errorMsg = "Failed to calculate EUI-64 address"

		slog.Error("EUI-64 calculation failed", "mac", mac, "prefix", prefix, "error", err)
		h.renderResult(c, http.StatusInternalServerError, interfaceID, fullIP, errorMsg)

		return
	}

	h.renderResult(c, http.StatusOK, interfaceID, fullIP, errorMsg)
}

// CalculateAdapter adapts the Calculate method to a gin.HandlerFunc.
func (h *Handler) CalculateAdapter() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.Calculate(c)
	}
}

// renderResult renders a calculation result with the specified status.
func (h *Handler) renderResult(c *gin.Context, status int, interfaceID, fullIP, errorMsg string) {
	c.Status(status) // Set status first

	if err := h.renderer.RenderResult(c, interfaceID, fullIP, errorMsg); err != nil {
		slog.Error("Failed to render result", "error", err)

		msg := errorMsg
		if msg == "" {
			msg = err.Error()
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})

		return
	}
}

// ginRequestContext adapts gin.Context to RequestContext.
type ginRequestContext struct {
	c *gin.Context
}

func (g *ginRequestContext) FormValue(key string) string {
	return g.c.PostForm(key)
}

func (g *ginRequestContext) GetContext() *gin.Context {
	return g.c
}
