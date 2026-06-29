package test

import (
	"net/http"
	"strings"
	"testing"

	"golazy.dev/lazytest"
	appinit "sample_app/init"
)

func TestApplicationRoutes(t *testing.T) {
	app := lazytest.New(t, appinit.App())

	app.Check(
		lazytest.Case{Name: "home", Method: http.MethodGet, Path: "/", Status: http.StatusOK, Contains: []string{"Hello", "helloworldservice", "data-controller=\"hello\""}, ContentType: "text/html"},
		lazytest.Case{Name: "public file", Method: http.MethodGet, Path: "/styles.css", Status: http.StatusOK, Contains: []string{"tailwindcss"}, ContentType: "text/css"},
		lazytest.Case{Name: "importmap", Method: http.MethodGet, Path: "/assets/importmap.json", Status: http.StatusOK, Contains: []string{"\"app.js\""}, ContentType: "application/json"},
		lazytest.Case{Name: "missing file", Method: http.MethodGet, Path: "/missing.txt", Status: http.StatusNotFound, Contains: []string{"404 page not found"}},
		lazytest.Case{Name: "liveness probe opt-in", Method: http.MethodGet, Path: "/livez", Status: http.StatusNotFound, Contains: []string{"404 page not found"}},
		lazytest.Case{Name: "readiness probe opt-in", Method: http.MethodGet, Path: "/readyz", Status: http.StatusNotFound, Contains: []string{"404 page not found"}},
		lazytest.Case{Name: "sitemap opt-in", Method: http.MethodGet, Path: "/sitemap.xml", Status: http.StatusNotFound, Contains: []string{"404 page not found"}},
		lazytest.Case{Name: "missing route below root", Method: http.MethodGet, Path: "/posts", Status: http.StatusNotFound, Contains: []string{"404 page not found"}},
		lazytest.Case{Name: "unsupported method", Method: http.MethodPost, Path: "/", Status: http.StatusMethodNotAllowed, Contains: []string{"Method Not Allowed"}, Allow: []string{http.MethodGet}},
		lazytest.Case{Name: "unsupported public method", Method: http.MethodPost, Path: "/styles.css", Status: http.StatusMethodNotAllowed, Contains: []string{"Method Not Allowed"}, Allow: []string{http.MethodGet}},
	)
}

func TestApplicationRouteTableOnlyExposesHome(t *testing.T) {
	app := appinit.App()
	if got, want := len(app.Router.Routes), 1; got != want {
		t.Fatalf("route count = %d, want %d: %#v", got, want, app.Router.Routes)
	}
	route := app.Router.Routes[0]
	if route.Method != http.MethodGet || route.Path != "/" || route.Action != "Index" || route.Controller != "home" {
		t.Fatalf("route = %#v, want GET / home#index", route)
	}
}

func TestApplicationUsesAssetPermalink(t *testing.T) {
	app := lazytest.New(t, appinit.App())

	home := app.Get("/")
	home.OK().
		Contains(`<html lang="en">`).
		Contains(`<script type="importmap">`).
		Contains(`import "app.js"`)

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
		"@hotwired/stimulus":              "/assets/lazyshaft/stimulus-",
		"@hotwired/turbo":                 "/assets/lazyshaft/turbo-",
		"app.js":                          "/assets/lazyshaft/app/app-",
		"controllers/hello_controller.js": "/assets/lazyshaft/app/controllers/hello_controller-",
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
