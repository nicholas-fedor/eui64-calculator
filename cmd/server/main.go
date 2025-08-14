// Package main provides the entry point for the EUI-64 calculator web server.
// It loads configuration from environment variables, sets up the Gin router with routes and middleware,
// and starts the HTTP server, handling errors by logging and exiting with a non-zero status.
package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
)

// Constants defining default configuration values and environment variable names.
const (
	defaultPort       = "8080"            // defaultPort is the default server port if PORT is unset.
	trustedProxiesEnv = "TRUSTED_PROXIES" // trustedProxiesEnv is the environment variable for trusted proxy IPs.
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
	ErrSetTrustedProxies = errors.New("failed to set trusted proxies")
	ErrSetupRouter       = errors.New("failed to setup router")
	ErrServerFailed      = errors.New("server failed")
)

// Config holds server configuration parameters.
type Config struct {
	Port           string   // Port is the server listening port (e.g., "8080").
	TrustedProxies []string // TrustedProxies lists IP addresses of trusted reverse proxies.
}

// LoadConfig loads server configuration from environment variables.
// It defaults to port ":8080" if PORT is unset and processes TRUSTED_PROXIES as a comma-separated list,
// trimming whitespace, logging warnings for empty entries, and filtering them out.
func LoadConfig() Config {
	config := Config{Port: ":" + defaultPort}
	if port := os.Getenv("PORT"); port != "" {
		config.Port = ":" + port
	}

	if proxies := os.Getenv(trustedProxiesEnv); proxies != "" {
		proxyList := strings.Split(proxies, ",")

		var validProxies []string

		for _, proxy := range proxyList {
			trimmedProxy := strings.TrimSpace(proxy)
			if trimmedProxy == "" {
				slog.WarnContext(context.Background(), "Empty proxy entry in TRUSTED_PROXIES")
			} else {
				validProxies = append(validProxies, trimmedProxy)
			}
		}

		config.TrustedProxies = validProxies
	}

	return config
}

// SetupRouter configures and returns a new Gin router with middleware and routes.
// It sets up the router in release mode, enables logging and recovery middleware, configures trusted proxies,
// and defines routes for the home page, EUI-64 calculation, and embedded file serving. Returns the router and any error.
func SetupRouter(config Config) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	gin.ForceConsoleColor() // Forces colored console output for logs, even in non-terminal environments.
	router.Use(gin.Logger(), gin.Recovery())

	if err := router.SetTrustedProxies(config.TrustedProxies); err != nil {
		return nil, errors.Join(ErrSetTrustedProxies, err)
	}

	// Create a sub-FS to serve files from the "static" subdirectory as if it were the root.
	subStatic, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, errors.Join(ErrSetupRouter, err)
	}

	handler := handlers.NewHandler(&eui64.DefaultCalculator{})
	router.GET("/", handler.Home)
	router.POST("/calculate", handler.Calculate)
	router.StaticFS("/static", http.FS(subStatic))

	return router, nil
}

// main initializes and runs the EUI-64 calculator web server.
// It loads configuration, sets up the router, and starts the server, logging errors and exiting with status 1 on failure.
func main() {
	config := LoadConfig()

	router, err := SetupRouter(config)
	if err != nil {
		slog.ErrorContext(context.Background(), ErrSetupRouter.Error(), "error", err)
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

	if err := router.Run(config.Port); err != nil {
		slog.ErrorContext(context.Background(), ErrServerFailed.Error(), "error", err)
		os.Exit(1)
	}
}
