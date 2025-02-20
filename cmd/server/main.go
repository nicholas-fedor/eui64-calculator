package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
)

// Constants defining default configuration values and environment variable names.
const (
	defaultPort       = ":8080"           // defaultPort is the default server port if PORT is unset.
	trustedProxiesEnv = "TRUSTED_PROXIES" // trustedProxiesEnv is the environment variable for trusted proxy IPs.
)

// Config holds server configuration parameters.
type Config struct {
	Port           string   // Port is the server listening port (e.g., ":8080").
	TrustedProxies []string // TrustedProxies lists IP addresses of trusted reverse proxies.
}

// LoadConfig loads server configuration from environment variables.
// It defaults to port ":8080" if PORT is unset and processes TRUSTED_PROXIES as a comma-separated list,
// trimming whitespace and logging warnings for empty entries. Returns the configuration and any error encountered.
func LoadConfig() (Config, error) {
	config := Config{Port: defaultPort}
	if port := os.Getenv("PORT"); port != "" {
		config.Port = ":" + port
	}
	if proxies := os.Getenv(trustedProxiesEnv); proxies != "" {
		config.TrustedProxies = strings.Split(proxies, ",")
		for i, proxy := range config.TrustedProxies {
			config.TrustedProxies[i] = strings.TrimSpace(proxy)
			if config.TrustedProxies[i] == "" {
				slog.Warn("Empty proxy entry in TRUSTED_PROXIES")
			}
		}
	}
	return config, nil
}

// SetupRouter configures and returns a new Gin router with middleware and routes.
// It sets up the router in release mode, enables logging and recovery middleware, configures trusted proxies,
// and defines routes for the home page, EUI-64 calculation, and static file serving. Returns the router and any error.
func SetupRouter(config Config) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	gin.ForceConsoleColor() // Forces colored console output for logs, even in non-terminal environments.
	r.Use(gin.Logger(), gin.Recovery())

	if err := r.SetTrustedProxies(config.TrustedProxies); err != nil {
		return nil, errors.Join(fmt.Errorf("failed to set trusted proxies"), err)
	}

	handler := handlers.NewHandler(&eui64.DefaultCalculator{})
	r.GET("/", handler.Home)
	r.POST("/calculate", handler.Calculate)
	r.Static("/static", "./static") // Serves static files from the project root's static directory.

	return r, nil
}

// main initializes and runs the EUI-64 calculator web server.
// It loads configuration, sets up the router, and starts the server, logging errors and exiting with status 1 on failure.
func main() {
	config, err := LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	router, err := SetupRouter(config)
	if err != nil {
		slog.Error("Failed to setup router", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting server", "port", config.Port)
	if err := router.Run(config.Port); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
