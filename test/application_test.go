package test

import (
	"net/http"
	"strings"
	"sync"
	"testing"

	"golazy.dev/lazytest"
	appinit "sample_app/init"
)

func TestApplicationRoutes(t *testing.T) {
	app := lazytest.New(t, appinit.App())

	app.Check(
		lazytest.Case{Name: "home", Method: http.MethodGet, Path: "/", Status: http.StatusOK, Contains: []string{"Hello, world!"}, ContentType: "text/html"},
		lazytest.Case{Name: "liveness probe", Method: http.MethodGet, Path: "/livez", Status: http.StatusOK, Contains: []string{"live"}, ContentType: "text/plain"},
		lazytest.Case{Name: "readiness probe", Method: http.MethodGet, Path: "/readyz", Status: http.StatusOK, Contains: []string{"ready"}, ContentType: "text/plain"},
		lazytest.Case{Name: "posts", Method: http.MethodGet, Path: "/posts", Status: http.StatusOK, Contains: []string{"Hello, GoLazy"}, ContentType: "text/html"},
		lazytest.Case{Name: "posts html suffix", Method: http.MethodGet, Path: "/posts.html", Status: http.StatusOK, Contains: []string{"Hello, GoLazy"}, ContentType: "text/html"},
		lazytest.Case{Name: "posts markdown", Method: http.MethodGet, Path: "/posts", Headers: map[string][]string{"Accept": {"text/markdown"}}, Status: http.StatusOK, Contains: []string{"- [Hello, GoLazy](/posts/hello-golazy)"}, ContentType: "text/markdown"},
		lazytest.Case{Name: "posts markdown suffix", Method: http.MethodGet, Path: "/posts.md", Headers: map[string][]string{"Accept": {"text/markdown"}}, Status: http.StatusOK, Contains: []string{"- [Hello, GoLazy](/posts/hello-golazy)"}, ContentType: "text/markdown"},
		lazytest.Case{Name: "post", Method: http.MethodGet, Path: "/posts/hello-golazy", Status: http.StatusOK, Contains: []string{"<strong>GoLazy</strong>"}, ContentType: "text/html"},
		lazytest.Case{Name: "post html suffix", Method: http.MethodGet, Path: "/posts/hello-golazy.html", Status: http.StatusOK, Contains: []string{"<strong>GoLazy</strong>"}, ContentType: "text/html"},
		lazytest.Case{Name: "post markdown", Method: http.MethodGet, Path: "/posts/hello-golazy", Headers: map[string][]string{"Accept": {"text/markdown"}}, Status: http.StatusOK, Contains: []string{"Welcome to **GoLazy**"}, ContentType: "text/markdown"},
		lazytest.Case{Name: "post markdown suffix", Method: http.MethodGet, Path: "/posts/hello-golazy.md", Headers: map[string][]string{"Accept": {"text/markdown"}}, Status: http.StatusOK, Contains: []string{"Welcome to **GoLazy**"}, ContentType: "text/markdown"},
		lazytest.Case{Name: "post word count helper", Method: http.MethodGet, Path: "/posts/hello-golazy", Status: http.StatusOK, Contains: []string{"25 words"}, ContentType: "text/html"},
		lazytest.Case{Name: "post read time helper", Method: http.MethodGet, Path: "/posts/hello-golazy", Status: http.StatusOK, Contains: []string{"1 min read"}, ContentType: "text/html"},
		lazytest.Case{Name: "missing post", Method: http.MethodGet, Path: "/posts/missing", Status: http.StatusNotFound, Contains: []string{"Not Found"}},
		lazytest.Case{Name: "public file", Method: http.MethodGet, Path: "/styles.css", Status: http.StatusOK, Contains: []string{"tailwindcss"}, ContentType: "text/css"},
		lazytest.Case{Name: "importmap", Method: http.MethodGet, Path: "/assets/importmap.json", Status: http.StatusOK, Contains: []string{"\"/js/app.js\""}, ContentType: "application/json"},
		lazytest.Case{Name: "missing file", Method: http.MethodGet, Path: "/missing.txt", Status: http.StatusNotFound, Contains: []string{"404 page not found"}},
		lazytest.Case{Name: "unsupported method", Method: http.MethodPost, Path: "/posts", Status: http.StatusMethodNotAllowed, Contains: []string{"Method Not Allowed"}, Allow: []string{http.MethodGet}},
		lazytest.Case{Name: "unsupported public method", Method: http.MethodPost, Path: "/styles.css", Status: http.StatusMethodNotAllowed, Contains: []string{"Method Not Allowed"}, Allow: []string{http.MethodGet}},
	)
}

func TestApplicationUsesAssetPermalink(t *testing.T) {
	app := lazytest.New(t, appinit.App())

	home := app.Get("/")
	home.OK().
		Contains(`<html lang="en" class="dark scheme-dark">`).
		Contains(`<script type="importmap">`).
		Contains(`import "/js/app.js"`)

	matches := home.Match(`href="(/styles-[a-f0-9]{12}\.css)"`)
	if len(matches) != 2 {
		t.Fatalf("home body does not contain fingerprinted stylesheet URL: %s", home.BodyString())
	}

	app.Get(matches[1]).
		OK().
		HeaderEquals("Cache-Control", "public, max-age=31536000, immutable")
}

func TestApplicationJavaScriptImportmap(t *testing.T) {
	app := lazytest.New(t, appinit.App())

	parsed := struct {
		Imports map[string]string `json:"imports"`
	}{}
	app.Get("/assets/importmap.json").OK().JSON(&parsed)

	for specifier, prefix := range map[string]string{
		"@hotwired/stimulus":                  "/assets/lazyshaft/stimulus-",
		"@hotwired/turbo":                     "/assets/lazyshaft/turbo-",
		"/js/app.js":                          "/assets/lazyshaft/app/app-",
		"/js/controllers/hello_controller.js": "/assets/lazyshaft/app/controllers/hello_controller-",
	} {
		assetPath, ok := parsed.Imports[specifier]
		if !ok {
			t.Fatalf("importmap did not contain %q", specifier)
		}
		if !strings.HasPrefix(assetPath, prefix) {
			t.Fatalf("importmap[%q] = %q, want prefix %q", specifier, assetPath, prefix)
		}
		app.Get(assetPath).OK().ContentType("text/javascript")
	}
}

func TestApplicationSessionsAndFlashes(t *testing.T) {
	t.Setenv("SECURE_COOKIE_KEY", "sample_app_session_test_key")

	browser := lazytest.New(t, appinit.App()).Client()

	response := browser.Get("/")
	response.OK().Contains("Visit count: 1")
	cookies := response.Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie after first visit")
	}

	browser.Get("/").OK().Contains("Visit count: 2")
	browser.Get("/flash").Status(http.StatusFound)
	browser.Get("/").
		OK().
		Contains("Visit count: 4").
		Contains("This is a sample flash message")
	browser.Get("/").OK().
		NotContains("This is a sample flash message")
}

func TestControllersHaveRequestLocalState(t *testing.T) {
	app := lazytest.New(t, appinit.App())

	var wait sync.WaitGroup
	for range 20 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			response := app.Get("/posts")
			if response.Result.StatusCode != http.StatusOK {
				t.Errorf("status = %d", response.Result.StatusCode)
			}
		}()
	}
	wait.Wait()
}
