package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/handlers/mocks"
	"github.com/stretchr/testify/require"
)

func TestCalculateHandlerInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		formData   url.Values
		validator  *mocks.Validator
		name       string
		wantBody   string
		wantStatus int
	}{
		{
			name: "Invalid MAC format",
			formData: url.Values{
				"mac":      {"invalid-mac"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			validator:  &mocks.Validator{MacErr: mocks.ErrInvalidMAC, PrefixErr: nil},
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
			wantStatus: http.StatusOK,
			wantBody:   "Please enter a valid IPv6 prefix (e.g., 2001:db8::)",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ginContext, responseRecorder := prepareCalcRequest(t, testCase.formData)
			handler, renderer := setupInvalidHandler(t, testCase.validator)
			ctx := mocks.NewRequestContext(ginContext) // Use constructor
			handler.Calculate(ctx)
			require.Equal(t, testCase.wantStatus, responseRecorder.Code)
			require.Contains(t, responseRecorder.Body.String(), testCase.wantBody)
			require.True(t, renderer.CalledResult, "RenderResult not called")
		})
	}
}
