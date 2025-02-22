package config

import (
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/nicholas-fedor/eui64-calculator/internal/utilities/config"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		port       string
		envProxies string
		wantConfig server.Config
	}{
		{
			name:       "Default config",
			port:       ":8080",
			envProxies: "",
			wantConfig: server.Config{
				Port:           ":8080",
				StaticDir:      "static",
				TrustedProxies: []string{},
			},
		},
		{
			name:       "With trusted proxies",
			port:       ":8080",
			envProxies: "192.168.1.1,10.0.0.1",
			wantConfig: server.Config{
				Port:           ":8080",
				StaticDir:      "static",
				TrustedProxies: []string{"192.168.1.1", "10.0.0.1"},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Removed t.Parallel() - using t.Setenv()
			if testCase.envProxies != "" {
				t.Setenv(config.TrustedProxiesEnv, testCase.envProxies)
			}

			got, err := config.LoadConfig(testCase.port)
			assertConfigResult(t, got, err, testCase.wantConfig)
		})
	}
}
