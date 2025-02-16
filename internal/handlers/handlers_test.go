package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter creates a Gin router for testing with the handlers set up.
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler()
	r.GET("/", handler.Home)
	r.POST("/calculate", handler.Calculate)
	return r
}

// TestHomeHandler tests the Home handler.
func TestHomeHandler(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Successful GET request",
			wantStatus: http.StatusOK,
			wantBody:   "EUI-64 Calculator", // Check for key content in the response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			req, _ := http.NewRequest("GET", "/", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Home() status code")
			assert.Contains(t, resp.Body.String(), tt.wantBody, "Home() response body")
		})
	}
}

// TestCalculateHandlerValid tests valid inputs to the Calculate handler.
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
			wantBody:   "0214:22ff:fe01:2345", // Check for interface ID in the response
		},
		{
			name: "Valid MAC and partial prefix",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:0db8"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "0214:22ff:fe01:2345", // Check for interface ID in the response
		},
		{
			name: "Valid MAC with no prefix",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {""},
			},
			wantStatus: http.StatusOK,
			wantBody:   "0214:22ff:fe01:2345", // Check for interface ID in the response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			req, _ := http.NewRequest("POST", "/calculate", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Calculate() status code")
			assert.Contains(t, resp.Body.String(), tt.wantBody, "Calculate() response body")
		})
	}
}

// TestCalculateHandlerInvalid tests invalid inputs to the Calculate handler.
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
			wantStatus: http.StatusOK,   // Assuming your app returns 200 with an error message
			wantBody:   "error-message", // Check for error message class in the response
		},
		{
			name: "MAC too short",
			formData: url.Values{
				"mac":      {"00-14-22-01-23"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "error-message",
		},
		{
			name: "Invalid prefix - too many hextets",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:0db8:85a3:0000:0000"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "error-message",
		},
		{
			name: "Invalid prefix - empty hextet",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001::85a3"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "error-message",
		},
		{
			name: "Invalid prefix - invalid hextet",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:invalid:85a3"},
			},
			wantStatus: http.StatusOK,
			wantBody:   "error-message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			req, _ := http.NewRequest("POST", "/calculate", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantStatus, resp.Code, "Calculate() status code")
			assert.Contains(t, resp.Body.String(), tt.wantBody, "Calculate() response body")
		})
	}
}
