package main

import (
	"net/http"
	"net/url"
	"testing"
)

func TestStaticRoute(t *testing.T) {
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
			name:       "GET /static/styles.css - Static file",
			method:     "GET",
			path:       "/static/styles.css",
			formData:   nil,
			wantBody:   "body",
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
