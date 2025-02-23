package main

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/app"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	// Backup and restore original newAppFunc
	origNewApp := newAppFunc
	defer func() { newAppFunc = origNewApp }()

	// Success case
	appInstance := &app.App{
		LoadConfig: func(_ string) (server.Config, error) {
			return server.Config{Port: ":8080", StaticDir: "static"}, nil
		},
		GinNew:    func(...gin.OptionFunc) *gin.Engine { return gin.New() },
		RunEngine: func(_ *gin.Engine, _ string) error { return nil },
		SetupRouter: func(_ server.Config, _, _ gin.HandlerFunc) (*gin.Engine, error) {
			return gin.New(), nil
		},
	}
	newAppFunc = func() *app.App { return appInstance }
	err := runFunc()
	require.NoError(t, err, "runFunc() should succeed")

	// Failure case
	appInstance.RunEngine = func(_ *gin.Engine, _ string) error { return errors.New("mock failure") }
	err = runFunc()
	require.Error(t, err, "runFunc() should fail")
}

func TestMainError(t *testing.T) {
	// Backup and restore original logFatalf
	origFatalf := logFatalf
	defer func() { logFatalf = origFatalf }()

	var exitCode int

	logFatalf = func(format string, v ...interface{}) {
		exitCode = 1
	}

	// Backup and restore original runFunc
	origRun := runFunc
	defer func() { runFunc = origRun }()

	runFunc = func() error { return errors.New("mock error") }

	main()
	require.Equal(t, 1, exitCode, "main() should call logFatalf on error")
}
