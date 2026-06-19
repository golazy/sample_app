package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	appinit "sample_app/init"
)

func TestApplicationRoutes(t *testing.T) {
	handler := application()

	tests := []struct {
		name        string
		method      string
		path        string
		accept      string
		status      int
		contains    string
		contentType string
	}{
		{name: "home", method: http.MethodGet, path: "/", status: http.StatusOK, contains: "Hello, world!", contentType: "text/html"},
		{name: "posts", method: http.MethodGet, path: "/posts", status: http.StatusOK, contains: "Hello, GoLazy", contentType: "text/html"},
		{name: "posts html suffix", method: http.MethodGet, path: "/posts.html", status: http.StatusOK, contains: "Hello, GoLazy", contentType: "text/html"},
		{name: "posts markdown", method: http.MethodGet, path: "/posts", accept: "text/markdown", status: http.StatusOK, contains: "- [Hello, GoLazy](/posts/hello-golazy)", contentType: "text/markdown"},
		{name: "posts markdown suffix", method: http.MethodGet, path: "/posts.md", status: http.StatusOK, contains: "- [Hello, GoLazy](/posts/hello-golazy)", contentType: "text/markdown"},
		{name: "post", method: http.MethodGet, path: "/posts/hello-golazy", status: http.StatusOK, contains: "<strong>GoLazy</strong>", contentType: "text/html"},
		{name: "post html suffix", method: http.MethodGet, path: "/posts/hello-golazy.html", status: http.StatusOK, contains: "<strong>GoLazy</strong>", contentType: "text/html"},
		{name: "post markdown", method: http.MethodGet, path: "/posts/hello-golazy", accept: "text/markdown", status: http.StatusOK, contains: "Welcome to **GoLazy**", contentType: "text/markdown"},
		{name: "post markdown suffix", method: http.MethodGet, path: "/posts/hello-golazy.md", status: http.StatusOK, contains: "Welcome to **GoLazy**", contentType: "text/markdown"},
		{name: "post word count helper", method: http.MethodGet, path: "/posts/hello-golazy", status: http.StatusOK, contains: "25 words", contentType: "text/html"},
		{name: "post read time helper", method: http.MethodGet, path: "/posts/hello-golazy", status: http.StatusOK, contains: "1 min read", contentType: "text/html"},
		{name: "missing post", method: http.MethodGet, path: "/posts/missing", status: http.StatusNotFound, contains: "Not Found"},
		{name: "public file", method: http.MethodGet, path: "/styles.css", status: http.StatusOK, contains: "tailwindcss", contentType: "text/css"},
		{name: "importmap", method: http.MethodGet, path: "/assets/importmap.json", status: http.StatusOK, contains: `"/js/app.js"`, contentType: "application/json"},
		{name: "missing file", method: http.MethodGet, path: "/missing.txt", status: http.StatusNotFound, contains: "404 page not found"},
		{name: "unsupported method", method: http.MethodPost, path: "/posts", status: http.StatusMethodNotAllowed, contains: "Method Not Allowed"},
		{name: "unsupported public method", method: http.MethodPost, path: "/styles.css", status: http.StatusMethodNotAllowed, contains: "Method Not Allowed"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			if test.accept != "" {
				request.Header.Set("Accept", test.accept)
			}
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

func TestApplicationUsesAssetPermalink(t *testing.T) {
	handler := application()

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if !strings.Contains(response.Body.String(), `<html lang="en" class="dark scheme-dark">`) {
		t.Fatalf("home body does not enable Tailwind dark mode by default: %s", response.Body.String())
	}

	matches := regexp.MustCompile(`href="(/styles-[a-f0-9]{12}\.css)"`).FindStringSubmatch(response.Body.String())
	if len(matches) != 2 {
		t.Fatalf("home body does not contain fingerprinted stylesheet URL: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `<script type="importmap">`) {
		t.Fatalf("home body does not contain inline importmap: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `import "/js/app.js"`) {
		t.Fatalf("home body does not import app JavaScript: %s", response.Body.String())
	}

	asset := httptest.NewRecorder()
	handler.ServeHTTP(asset, httptest.NewRequest(http.MethodGet, matches[1], nil))
	if asset.Code != http.StatusOK {
		t.Fatalf("asset status = %d, want %d", asset.Code, http.StatusOK)
	}
	if got := asset.Header().Get("Cache-Control"); got != "public, max-age=31536000, immutable" {
		t.Fatalf("asset Cache-Control = %q, want immutable cache policy", got)
	}
}

func TestApplicationJavaScriptImportmap(t *testing.T) {
	handler := application()

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/assets/importmap.json", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var parsed struct {
		Imports map[string]string `json:"imports"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("parse importmap: %v", err)
	}

	for specifier, prefix := range map[string]string{
		"@hotwired/stimulus":                  "/assets/lazyshaft/stimulus-",
		"@hotwired/turbo":                     "/assets/lazyshaft/turbo-",
		"/js/app.js":                          "/assets/lazyshaft/app/app-",
		"/js/controllers/hello_controller.js": "/assets/lazyshaft/app/controllers/hello_controller-",
	} {
		assetPath := parsed.Imports[specifier]
		if !strings.HasPrefix(assetPath, prefix) {
			t.Fatalf("importmap[%q] = %q, want prefix %q", specifier, assetPath, prefix)
		}

		asset := httptest.NewRecorder()
		handler.ServeHTTP(asset, httptest.NewRequest(http.MethodGet, assetPath, nil))
		if asset.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", assetPath, asset.Code, http.StatusOK)
		}
		if !strings.Contains(asset.Header().Get("Content-Type"), "text/javascript") {
			t.Fatalf("%s Content-Type = %q, want JavaScript", assetPath, asset.Header().Get("Content-Type"))
		}
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
