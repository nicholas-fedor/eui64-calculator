package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nicholas-fedor/eui64-calculator/internal/handlers"
	"github.com/nicholas-fedor/eui64-calculator/internal/handlers/mocks"
	"github.com/stretchr/testify/require"
)

func TestCalculateHandlerValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		formData   url.Values
		name       string
		wantBody   string
		wantStatus int
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ginContext, responseRecorder := prepareCalcRequest(t, testCase.formData)
			calc := &mocks.Calculator{
				InterfaceID: "0214:22ff:fe01:2345",
				FullIP:      "2001:db8::214:22ff:fe01:2345",
				Err:         nil,
			}
			renderer := &mocks.Renderer{
				HomeErr:      nil,
				ResultErr:    nil,
				CalledHome:   false,
				CalledResult: false,
			}
			handler := handlers.NewHandler(calc, &mocks.Validator{MacErr: nil, PrefixErr: nil}, renderer)
			ctx := mocks.NewRequestContext(ginContext) // Use constructor
			handler.Calculate(ctx)
			require.Equal(t, testCase.wantStatus, responseRecorder.Code)
			require.Contains(t, responseRecorder.Body.String(), testCase.wantBody)
			require.True(t, renderer.CalledResult, "RenderResult not called")
		})
	}
}
