// Package handlers provides HTTP request handlers for the EUI-64 calculator
// application using the Fiber framework. It defines the Handler struct with
// dependency injection for the EUI-64 calculator, and includes handlers for
// rendering the home page, processing calculation requests with validation,
// and rendering results or errors.
package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"github.com/nicholas-fedor/eui64-calculator/internal/ui"
	"github.com/nicholas-fedor/eui64-calculator/internal/validators"
)

// Calculator defines the interface for EUI-64 calculation logic used by handlers.
type Calculator interface {
	// CalculateEUI64 computes the EUI-64 interface ID and full IPv6 address
	// from a MAC address and prefix.
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

// Calculate handles POST requests to compute an EUI-64 address from form data.
// It validates the MAC address and IPv6 prefix from the request, computes
// the EUI-64 interface ID and full IPv6 address, and renders the result.
// Errors during validation or calculation are logged and displayed to the user.
func (h *Handler) Calculate(c fiber.Ctx) error {
	mac := c.FormValue("mac")
	prefix := c.FormValue("ip-start")
	data := ui.ResultData{}

	if err := validators.ValidateMAC(mac); err != nil {
		data.Error = "Please enter a valid MAC address (e.g., 00-14-22-01-23-45)"

		slog.WarnContext(
			c.Context(),
			"MAC validation failed",
			"mac", mac,
			"error", err,
		)

		return h.renderResult(c, data)
	}

	if err := validators.ValidateIPv6Prefix(prefix); err != nil {
		data.Error = "Please enter a valid IPv6 prefix (e.g., 2001:db8::)"

		slog.WarnContext(
			c.Context(),
			"Prefix validation failed",
			"prefix",
			prefix,
			"error",
			err,
		)

		return h.renderResult(c, data)
	}

	interfaceID, fullIP, err := h.calc.CalculateEUI64(mac, prefix)
	data.InterfaceID = interfaceID
	data.FullIP = fullIP

	if err != nil {
		data.Error = "Failed to calculate EUI-64 address"

		slog.ErrorContext(
			c.Context(),
			"EUI-64 calculation failed",
			"mac",
			mac,
			"prefix",
			prefix,
			"error",
			err,
		)
	}

	return h.renderResult(c, data)
}

// Home handles GET requests to the root path, rendering the home page.
// It serves the initial form for entering MAC and IPv6 prefix values,
// aborting with a 500 status on render failure.
func (h *Handler) Home(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")

	err := ui.Home().Render(
		c.Context(),
		c.Response().BodyWriter(),
	)
	if err != nil {
		slog.ErrorContext(
			c.Context(),
			"Failed to render home page",
			"error", err,
		)

		return fmt.Errorf("send internal server error: %w", c.SendStatus(http.StatusInternalServerError))
	}

	return nil
}

// renderResult renders the calculation result to the HTTP response.
// It uses the provided ResultData to display either the computed EUI-64 address
// or an error message, returning a 500 status if rendering fails.
func (h *Handler) renderResult(c fiber.Ctx, data ui.ResultData) error {
	c.Set("Content-Type", "text/html; charset=utf-8")

	err := ui.Result(data).Render(
		c.Context(),
		c.Response().BodyWriter(),
	)
	if err != nil {
		slog.ErrorContext(
			c.Context(),
			"Failed to render result",
			"error", err,
		)

		return fmt.Errorf("send internal server error: %w", c.SendStatus(http.StatusInternalServerError))
	}

	return nil
}
