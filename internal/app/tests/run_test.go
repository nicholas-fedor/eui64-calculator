package app

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/app"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
)

var (
	ErrConfigLoadFailed = errors.New("config load failed")
	ErrRouterRunFailed  = errors.New("router run failed")
	ErrSetupFailed      = errors.New("setup failed")
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		configPort string
		configErr  error
		setupErr   error
		runErr     error
		wantErr    bool
	}{
		{
			name:       "Successful run",
			configPort: ":0",
			configErr:  nil,
			setupErr:   nil,
			runErr:     nil,
			wantErr:    false,
		},
		{
			name:       "Config load error",
			configPort: ":0",
			configErr:  ErrConfigLoadFailed,
			setupErr:   nil,
			runErr:     nil,
			wantErr:    true,
		},
		{
			name:       "Router setup error",
			configPort: ":0",
			configErr:  nil,
			setupErr:   ErrSetupFailed,
			runErr:     nil,
			wantErr:    true,
		},
		{
			name:       "Router run error",
			configPort: ":0",
			configErr:  nil,
			setupErr:   nil,
			runErr:     ErrRouterRunFailed,
			wantErr:    true,
		},
		{
			name:       "Default NewApp run",
			configPort: ":8082",
			configErr:  nil,
			setupErr:   nil,
			runErr:     ErrRouterRunFailed,
			wantErr:    true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			runTestCase(t, testCase)
		})
	}
}

// runTestCase executes a single TestRun case with setup and assertions.
func runTestCase(t *testing.T, testCase struct {
	name       string
	configPort string
	configErr  error
	setupErr   error
	runErr     error
	wantErr    bool
}) {
	t.Helper()

	appInstance := app.NewApp()
	origLoadConfig := appInstance.LoadConfig
	origSetupRouter := appInstance.SetupRouter
	origRunEngine := appInstance.RunEngine

	defer func() {
		appInstance.LoadConfig = origLoadConfig
		appInstance.SetupRouter = origSetupRouter
		appInstance.RunEngine = origRunEngine
	}()

	appInstance.LoadConfig = func(_ string) (server.Config, error) {
		if testCase.configErr != nil {
			return server.Config{}, testCase.configErr
		}

		return server.Config{
			Port:           testCase.configPort,
			StaticDir:      "/tmp/static",
			TrustedProxies: []string{"127.0.0.1"},
		}, nil
	}

	appInstance.SetupRouter = func(_ server.Config, _, _ gin.HandlerFunc) (*gin.Engine, error) {
		if testCase.setupErr != nil {
			return nil, testCase.setupErr
		}

		return gin.New(), nil
	}

	appInstance.RunEngine = func(_ *gin.Engine, _ string) error {
		return testCase.runErr
	}

	err := appInstance.Run()
	assertRunError(t, err, testCase.wantErr, testCase.runErr)
}
