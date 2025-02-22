package tests

import (
	"net/http"
	"net/url"
	"testing"
)

func TestHomeRoute(t *testing.T) {
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
			name:       "GET / - Home page",
			method:     "GET",
			path:       "/",
			formData:   nil,
			wantBody:   "EUI-64 Calculator",
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET /nonexistent - Not found",
			method:     "GET",
			path:       "/nonexistent",
			formData:   nil,
			wantBody:   "404 page not found",
			wantStatus: http.StatusNotFound,
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
