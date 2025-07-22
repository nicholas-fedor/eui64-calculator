// Package main provides the entry point for the EUI-64 calculator web server.
// It loads configuration from environment variables, sets up the Gin router with routes and middleware,
// and starts the HTTP server, handling errors by logging and exiting with a non-zero status.
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
)

// Constants defining default configuration values and environment variable names.
const (
	defaultPort       = "8080"            // defaultPort is the default server port if PORT is unset.
	trustedProxiesEnv = "TRUSTED_PROXIES" // trustedProxiesEnv is the environment variable for trusted proxy IPs.
	staticDirEnv      = "STATIC_DIR"      // staticDirEnv is the environment variable for the static directory.
	defaultStaticDir  = "static"          // defaultStaticDir is the default static directory relative to the project root.
)

// Define static error variables.
var (
	ErrSetTrustedProxies = errors.New("failed to set trusted proxies")
	ErrLoadConfig        = errors.New("failed to load config")
	ErrSetupRouter       = errors.New("failed to setup router")
	ErrServerFailed      = errors.New("server failed")
)

// Config holds server configuration parameters.
type Config struct {
	Port           string   // Port is the server listening port (e.g., "8080").
	TrustedProxies []string // TrustedProxies lists IP addresses of trusted reverse proxies.
	StaticDir      string   // StaticDir is the directory containing static files.
}

// LoadConfig loads server configuration from environment variables.
// It defaults to port ":8080" if PORT is unset and processes TRUSTED_PROXIES as a comma-separated list,
// trimming whitespace, logging warnings for empty entries, and filtering them out.
func LoadConfig() Config {
	config := Config{Port: ":" + defaultPort, StaticDir: defaultStaticDir}
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

	if staticDir := os.Getenv(staticDirEnv); staticDir != "" {
		config.StaticDir = staticDir
	} else {
		exePath, err := os.Executable()
		if err != nil {
			slog.WarnContext(context.Background(), "Failed to determine executable path, using default static dir", "error", err)
		} else {
			config.StaticDir = filepath.Join(filepath.Dir(exePath), defaultStaticDir)
		}
	}

	return config
}

// SetupRouter configures and returns a new Gin router with middleware and routes.
// It sets up the router in release mode, enables logging and recovery middleware, configures trusted proxies,
// and defines routes for the home page, EUI-64 calculation, and static file serving. Returns the router and any error.
func SetupRouter(config Config) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	gin.ForceConsoleColor() // Forces colored console output for logs, even in non-terminal environments.
	router.Use(gin.Logger(), gin.Recovery())

	if err := router.SetTrustedProxies(config.TrustedProxies); err != nil {
		return nil, errors.Join(ErrSetTrustedProxies, err)
	}

	handler := handlers.NewHandler(&eui64.DefaultCalculator{})
	router.GET("/", handler.Home)
	router.POST("/calculate", handler.Calculate)
	router.Static("/static", config.StaticDir) // Use configurable static directory

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
		"port",
		config.Port,
		"static_dir",
		config.StaticDir,
	)

	if err := router.Run(config.Port); err != nil {
		slog.ErrorContext(context.Background(), ErrServerFailed.Error(), "error", err)
		os.Exit(1)
	}
}
