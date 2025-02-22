package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestCalculateRouteValid(t *testing.T) {
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
			name:   "POST /calculate - Valid MAC and full prefix",
			method: "POST",
			path:   "/calculate",
			formData: url.Values{
				"mac":      {"00-14-22-01-23-45"},
				"ip-start": {"2001:0db8:85a3:0000"},
			},
			wantBody:   "0214:22ff:fe01:2345",
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
