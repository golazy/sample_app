package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	appinit "sample_app/app/init"
)

func TestApplicationRoutes(t *testing.T) {
	handler := application()

	tests := []struct {
		name        string
		method      string
		path        string
		status      int
		contains    string
		contentType string
	}{
		{name: "home", method: http.MethodGet, path: "/", status: http.StatusOK, contains: "Hello, world!", contentType: "text/html"},
		{name: "posts", method: http.MethodGet, path: "/posts", status: http.StatusOK, contains: "Hello, GoLazy", contentType: "text/html"},
		{name: "post", method: http.MethodGet, path: "/posts/hello-golazy", status: http.StatusOK, contains: "<strong>GoLazy</strong>", contentType: "text/html"},
		{name: "missing post", method: http.MethodGet, path: "/posts/missing", status: http.StatusNotFound, contains: "Not Found"},
		{name: "public file", method: http.MethodGet, path: "/styles.css", status: http.StatusOK, contains: "color-scheme", contentType: "text/css"},
		{name: "missing file", method: http.MethodGet, path: "/missing.txt", status: http.StatusNotFound, contains: "404 page not found"},
		{name: "unsupported method", method: http.MethodPost, path: "/posts", status: http.StatusMethodNotAllowed, contains: "Method Not Allowed"},
		{name: "unsupported public method", method: http.MethodPost, path: "/styles.css", status: http.StatusMethodNotAllowed, contains: "Method Not Allowed"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			if response.Code != test.status {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, test.status, response.Body.String())
			}
			if !strings.Contains(response.Body.String(), test.contains) {
				t.Fatalf("body %q does not contain %q", response.Body.String(), test.contains)
			}
			if test.contentType != "" && !strings.Contains(response.Header().Get("Content-Type"), test.contentType) {
				t.Fatalf("Content-Type = %q, want %q", response.Header().Get("Content-Type"), test.contentType)
			}
			if test.status == http.StatusMethodNotAllowed && !strings.Contains(response.Header().Get("Allow"), http.MethodGet) {
				t.Fatalf("Allow = %q, want it to contain %q", response.Header().Get("Allow"), http.MethodGet)
			}
		})
	}
}

func TestControllersHaveRequestLocalState(t *testing.T) {
	handler := application()

	var wait sync.WaitGroup
	for range 20 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/posts", nil))
			if response.Code != http.StatusOK {
				t.Errorf("status = %d", response.Code)
			}
		}()
	}
	wait.Wait()
}

func application() http.Handler {
	return appinit.App()
}
