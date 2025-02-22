package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/stretchr/testify/require"
)

// assertConfigResult checks the LoadConfig result against expectations.
func assertConfigResult(t *testing.T, got server.Config, err error, wantConfig server.Config) {
	t.Helper()
	require.NoError(t, err)

	// Adjust StaticDir expectation based on os.Executable() behavior
	if got.StaticDir != "static" {
		exePath, _ := os.Executable()
		require.Equal(t, filepath.Join(filepath.Dir(exePath), "static"), got.StaticDir)
		got.StaticDir = "static" // Normalize for comparison
	}

	require.Equal(t, wantConfig, got)
}
