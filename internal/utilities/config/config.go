package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicholas-fedor/eui64-calculator/internal/server"
)

const (
	TrustedProxiesEnv = "TRUSTED_PROXIES"
	DefaultStaticDir  = "static"
)

// LoadConfig loads the server configuration from environment variables or defaults.
func LoadConfig(port string) (server.Config, error) {
	config := server.Config{
		Port:           port,
		StaticDir:      "",
		TrustedProxies: []string{},
	}

	if trustedProxies, ok := os.LookupEnv(TrustedProxiesEnv); ok {
		config.TrustedProxies = strings.Split(trustedProxies, ",")
	}

	exePath, err := os.Executable()
	if err != nil {
		log.Printf("WARN failed to get executable path: %v", err)

		config.StaticDir = DefaultStaticDir

		return config, nil
	}

	config.StaticDir = filepath.Join(filepath.Dir(exePath), DefaultStaticDir)

	return config, nil
}
