// Package main provides the entry point for the EUI-64 calculator web server.
// It loads configuration from environment variables, sets up the Fiber app with
// routes and middleware, and starts the HTTP server, handling errors by logging
// and exiting with a non-zero status.
package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
)

// Config holds server configuration parameters.
type Config struct {
	// Port is the server listening port (e.g., "8080").
	Port string
	// TrustedProxies lists IP addresses of trusted reverse proxies.
	TrustedProxies []string
}

// Constants defining default configuration values and environment variable names.
const (
	// defaultPort is the default server port if PORT is unset.
	defaultPort = "8080"
	// trustedProxiesEnv is the environment variable for trusted proxy IPs.
	trustedProxiesEnv = "TRUSTED_PROXIES"
)

//go:embed static/*
var staticFS embed.FS // Embeds all files in cmd/server/static/

// Build information injected by GoReleaser.
var (
	version string // Commit tag without "v" prefix
	commit  string // Commit SHA digest
	date    string // Commit date
)

// Define static error variables.
var (
	ErrSetupRouter  = errors.New("failed to setup router")
	ErrServerFailed = errors.New("server failed")
)

// LoadConfig loads server configuration from environment variables.
// It defaults to port ":8080" if PORT is unset and processes TRUSTED_PROXIES
// as a comma-separated list, trimming whitespace, logging warnings for empty
// entries, and filtering them out.
func LoadConfig() Config {
	config := Config{
		Port:           ":" + defaultPort,
		TrustedProxies: nil,
	}
	if port := os.Getenv("PORT"); port != "" {
		config.Port = ":" + port
	}

	if proxies := os.Getenv(trustedProxiesEnv); proxies != "" {
		proxyList := strings.Split(proxies, ",")

		var validProxies []string

		for _, proxy := range proxyList {
			trimmedProxy := strings.TrimSpace(proxy)
			if trimmedProxy == "" {
				slog.WarnContext(
					context.Background(),
					"Empty proxy entry in TRUSTED_PROXIES",
				)
			} else {
				validProxies = append(validProxies, trimmedProxy)
			}
		}

		config.TrustedProxies = validProxies
	}

	return config
}

// SetupRouter configures and returns a new Fiber app with middleware and routes.
// It sets up logging and recovery middleware, configures trusted proxies,
// and defines routes for the home page, EUI-64 calculation, and embedded file serving.
// Returns the app and any error.
func SetupRouter(config Config) (*fiber.App, error) {
	fiberCfg := fiber.Config{}

	if len(config.TrustedProxies) > 0 {
		fiberCfg.TrustProxy = true
		fiberCfg.TrustProxyConfig = fiber.TrustProxyConfig{
			Proxies:    config.TrustedProxies,
			LinkLocal:  false,
			Loopback:   false,
			Private:    false,
			UnixSocket: false,
		}
		fiberCfg.ProxyHeader = fiber.HeaderXForwardedFor
		fiberCfg.EnableIPValidation = true
	}

	app := fiber.New(fiberCfg)

	app.Use(recover.New(), logger.New())

	// Create a sub-FS to serve files from the "static" subdirectory as if it were the root.
	subStatic, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, errors.Join(ErrSetupRouter, err)
	}

	app.Use("/static", static.New("", static.Config{
		FS:              subStatic,
		Next:            nil,
		ModifyResponse:  nil,
		NotFoundHandler: nil,
		IndexNames:      []string{"index.html"},
		CacheDuration:   0,
		MaxAge:          0,
		Compress:        false,
		ByteRange:       false,
		Browse:          false,
		Download:        false,
	}))

	handler := handlers.NewHandler(&eui64.DefaultCalculator{})
	app.Get("/", handler.Home)
	app.Post("/calculate", handler.Calculate)

	return app, nil
}

// main initializes and runs the EUI-64 calculator web server.
// It loads configuration, sets up the app, and starts the server,
// logging errors and exiting with status 1 on failure.
func main() {
	config := LoadConfig()

	app, err := SetupRouter(config)
	if err != nil {
		slog.ErrorContext(
			context.Background(),
			ErrSetupRouter.Error(),
			"error", err,
		)
		os.Exit(1)
	}

	slog.InfoContext(
		context.Background(),
		"Starting server",
		"port", config.Port,
		"version", version,
		"commit", commit,
		"build_date", date,
	)

	if err := app.Listen(config.Port); err != nil {
		slog.ErrorContext(
			context.Background(),
			ErrServerFailed.Error(),
			"error", err,
		)
		os.Exit(1)
	}
}
