package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
	"github.com/nicholas-fedor/eui64-calculator/internal/server"
	"github.com/nicholas-fedor/eui64-calculator/internal/utilities/config"
	"github.com/nicholas-fedor/eui64-calculator/internal/validators"
	"github.com/stretchr/testify/require"
)

func TestRouterSetupIntegration(t *testing.T) {
	t.Parallel()

	config, err := config.LoadConfig(":0")
	require.NoError(t, err)

	calculator := &eui64.DefaultCalculator{}
	validator := &validators.CombinedValidator{}
	handler := handlers.NewHandler(calculator, validator, &server.UIRenderer{})

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	err = router.SetTrustedProxies(config.TrustedProxies)
	require.NoError(t, err)

	router.GET("/", handler.HomeAdapter())
	router.POST("/calculate", handler.CalculateAdapter())
	router.Static("/static", config.StaticDir)

	srv := httptest.NewServer(router)
	defer srv.Close()

	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL+"/", nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	form := url.Values{
		"mac":      {"00-14-22-01-23-45"},
		"ip-start": {"2001:0db8:85a3:0000"},
	}
	req, err = http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		srv.URL+"/calculate",
		strings.NewReader(form.Encode()),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
