package tests

import (
	"net/http"
	"net/url"
	"testing"
)

func TestCalculateRouteInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		formData   url.Values
		wantBody   string
		wantStatus int
	}{
		{
			name:   "POST /calculate - Invalid MAC",
			method: "POST",
			path:   "/calculate",
			formData: url.Values{
				"mac":      {"invalid-mac"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			wantBody:   "error-message",
			wantStatus: http.StatusOK,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			router := setupRouter(t)
			req := prepareRouteRequest(t, testCase.method, testCase.path, testCase.formData)
			testRoute(t, router, req, testCase.wantStatus, testCase.wantBody)
		})
	}
}
