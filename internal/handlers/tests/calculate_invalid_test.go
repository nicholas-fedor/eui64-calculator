package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/handlers/mocks"
	"github.com/stretchr/testify/require"
)

func TestCalculateHandlerInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		formData   url.Values
		validator  *mocks.Validator
		calculator *mocks.Calculator
		renderer   *mocks.MockRenderer // Use MockRenderer
		wantStatus int
		wantBody   string
	}{
		{
			name: "Invalid MAC format",
			formData: url.Values{
				"mac":      {"invalid-mac"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			validator:  &mocks.Validator{MacErr: mocks.ErrInvalidMAC, PrefixErr: nil},
			calculator: &mocks.Calculator{},
			renderer:   &mocks.MockRenderer{},
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid MAC address (e.g., 00-14-22-01-23-45)",
		},
		{
			name: "Invalid prefix",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001::85a3"},
			},
			validator:  &mocks.Validator{MacErr: nil, PrefixErr: mocks.ErrInvalidPrefix},
			calculator: &mocks.Calculator{},
			renderer:   &mocks.MockRenderer{},
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid IPv6 prefix (e.g., 2001:db8::)",
		},
		{
			name: "Renderer failure",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			validator:  &mocks.Validator{MacErr: nil, PrefixErr: nil},
			calculator: &mocks.Calculator{},
			renderer:   &mocks.MockRenderer{ResultErr: errors.New("render failed")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "render failed",
		},
		{
			name: "Calculation failure",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			validator:  &mocks.Validator{MacErr: nil, PrefixErr: nil},
			calculator: &mocks.Calculator{Err: errors.New("calculation failed")},
			renderer:   &mocks.MockRenderer{},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "Failed to calculate EUI-64 address",
		},
		{
			name: "Empty MAC",
			formData: url.Values{
				"mac":      {""},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			validator:  &mocks.Validator{MacErr: errors.New("MAC required"), PrefixErr: nil},
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid MAC address",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ginContext, responseRecorder := prepareCalcRequest(t, testCase.formData)
			handler, renderer := setupInvalidHandler(t, testCase.validator, testCase.calculator, testCase.renderer)
			handler.Calculate(ginContext) // Pass *gin.Context directly
			require.Equal(t, testCase.wantStatus, responseRecorder.Code)
			require.Contains(t, responseRecorder.Body.String(), testCase.wantBody)
			require.True(t, renderer.CalledResult, "RenderResult not called")
		})
	}
}
