package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/stretchr/testify/assert"
)

// setupRouter creates a Gin router for testing handler functions.
// It configures the router in test mode with the default EUI-64 calculator, setting up routes for home and calculate endpoints.
func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler(&eui64.DefaultCalculator{})
	r.GET("/", handler.Home)
	r.POST("/calculate", handler.Calculate)
	return r
}

// TestHomeHandler tests the Home handler’s response to GET requests.
// It verifies that the handler renders the home page with a 200 status and includes expected content.
func TestHomeHandler(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Successful GET request",
			wantStatus: http.StatusOK,
			wantBody:   "EUI-64 Calculator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(t)
			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantStatus, resp.Code)
			assert.Contains(t, resp.Body.String(), tt.wantBody)
		})
	}
}

// TestCalculateHandlerValid tests the Calculate handler with valid form inputs.
// It ensures the handler processes valid MAC addresses and IPv6 prefixes, returning a 200 status
// and the correct EUI-64 interface ID or full IPv6 address in the response.
func TestCalculateHandlerValid(t *testing.T) {
	tests := []struct {
		name       string
		formData   url.Values
		wantStatus int
		wantBody   string
	}{
		{
			name: "Valid MAC and full prefix",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "0214:22ff:fe01:2345",
		},
		{
			name: "Valid MAC with valid prefix",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:db8::"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "2001:db8::214:22ff:fe01:2345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(t)
			req, _ := http.NewRequest("POST", "/calculate", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantStatus, resp.Code)
			assert.Contains(t, resp.Body.String(), tt.wantBody)
		})
	}
}

// TestCalculateHandlerInvalid tests the Calculate handler with invalid form inputs.
// It verifies that the handler returns a 200 status with appropriate error messages
// for malformed MAC addresses and IPv6 prefixes, ensuring proper validation feedback.
func TestCalculateHandlerInvalid(t *testing.T) {
	tests := []struct {
		name       string
		formData   url.Values
		wantStatus int
		wantBody   string
	}{
		{
			name: "Invalid MAC format",
			formData: url.Values{
				"mac":      {"invalid-mac"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid MAC address (e.g., 00-14-22-01-23-45)",
		},
		{
			name: "MAC too short",
			formData: url.Values{
				"mac":      {"00-14-22-01-23"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid MAC address (e.g., 00-14-22-01-23-45)",
		},
		{
			name: "Invalid prefix - too many hextets",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:0db8:85a3:0000:0000"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid IPv6 prefix (e.g., 2001:db8::)",
		},
		{
			name: "Invalid prefix - empty hextet",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001::85a3"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid IPv6 prefix (e.g., 2001:db8::)",
		},
		{
			name: "Invalid prefix - invalid hextet",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:invalid:85a3"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid IPv6 prefix (e.g., 2001:db8::)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(t)
			req, _ := http.NewRequest("POST", "/calculate", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.wantStatus, resp.Code)
			assert.Contains(t, resp.Body.String(), tt.wantBody)
		})
	}
}
