package app

import (
	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/nicholas-fedor/eui64-calculator/internal/utilities/config"
)

const DefaultPort = ":8080"

type LoadConfigFunc func(string) (server.Config, error)
type GinNewFunc func(...gin.OptionFunc) *gin.Engine
type RunFunc func(*gin.Engine, string) error
type SetupRouterFunc func(server.Config, gin.HandlerFunc, gin.HandlerFunc) (*gin.Engine, error)

// App encapsulates the application’s runtime dependencies.
type App struct {
	LoadConfig  LoadConfigFunc
	GinNew      GinNewFunc
	RunEngine   RunFunc
	SetupRouter SetupRouterFunc
}

// NewApp creates a new App instance with default dependencies.
func NewApp() *App {
	return &App{
		LoadConfig:  config.LoadConfig,
		GinNew:      gin.New,
		RunEngine:   func(e *gin.Engine, addr string) error { return e.Run(addr) },
		SetupRouter: server.SetupRouter,
	}
}

// Run starts the application with the configured dependencies.
func (a *App) Run() error {
	config, err := a.LoadConfig(DefaultPort)
	if err != nil {
		return err
	}

	router, err := a.SetupRouter(config, nil, nil)
	if err != nil {
		return err
	}

	return a.RunEngine(router, config.Port)
}
