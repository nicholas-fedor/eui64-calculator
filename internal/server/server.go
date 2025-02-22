package server

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
	"github.com/nicholas-fedor/eui64-calculator/internal/validators"
	"github.com/nicholas-fedor/eui64-calculator/ui"
)

type Config struct {
	Port           string
	StaticDir      string
	TrustedProxies []string
}

var (
	ErrSetTrustedProxies = errors.New("failed to set trusted proxies")
	ErrRenderFailed      = errors.New("failed to render UI component")
)

type UIRenderer struct{}

func (r *UIRenderer) RenderHome(ctx *gin.Context) error {
	if err := ui.Home().Render(ctx.Request.Context(), ctx.Writer); err != nil {
		return fmt.Errorf("%w: %w", ErrRenderFailed, err)
	}

	return nil
}

func (r *UIRenderer) RenderResult(ctx *gin.Context, interfaceID, fullIP, errorMsg string) error {
	data := ui.ResultData{
		InterfaceID: interfaceID,
		FullIP:      fullIP,
		Error:       errorMsg,
	}
	if err := ui.Result(data).Render(ctx.Request.Context(), ctx.Writer); err != nil {
		return fmt.Errorf("%w: %w", ErrRenderFailed, err)
	}

	return nil
}

func SetupRouter(config Config, homeHandler, calculateHandler gin.HandlerFunc) (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	if err := router.SetTrustedProxies(config.TrustedProxies); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSetTrustedProxies, err)
	}

	// Use provided handlers if not nil, otherwise use defaults
	if homeHandler == nil {
		calculator := &eui64.DefaultCalculator{}
		validator := &validators.CombinedValidator{}
		handler := handlers.NewHandler(calculator, validator, &UIRenderer{})
		homeHandler = handler.HomeAdapter()
	}

	if calculateHandler == nil {
		calculator := &eui64.DefaultCalculator{}
		validator := &validators.CombinedValidator{}
		handler := handlers.NewHandler(calculator, validator, &UIRenderer{})
		calculateHandler = handler.CalculateAdapter()
	}

	router.GET("/", homeHandler)
	router.POST("/calculate", calculateHandler)
	router.Static("/static", config.StaticDir)

	return router, nil
}
