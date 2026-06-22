package appinit

import (
	"os"
	"time"

	"golazy.dev/lazyapp"
	"golazy.dev/lazycontrolplane"
	"golazy.dev/lazyseo"
	"golazy.dev/lazysession"
	_ "golazy.dev/lazyview/gotmpl"
	"sample_app/app"
	"sample_app/app/helpers"
)

func App() *lazyapp.App {
	return lazyapp.New(lazyapp.Config{
		Name:         "sample_app",
		Drawer:       Draw,
		Public:       app.Public,
		Views:        app.Views,
		Context:      Context,
		Helpers:      lazyapp.Helpers{helpers.RegisterHelpers()},
		ControlPlane: lazycontrolplane.Config{},
		SEO: []lazyseo.Option{
			lazyseo.SiteName("GoLazy"),
			lazyseo.Description("A small GoLazy sample application with embedded templates, assets, and Markdown posts."),
			lazyseo.Language("en"),
			lazyseo.Type("website"),
			lazyseo.Locale("en_US"),
			lazyseo.TwitterCardType("summary"),
		},
		Sitemap: lazyapp.SitemapConfig{
			URLs: sampleSitemapURLs(),
		},
		Sessions: lazysession.Config{
			Key: os.Getenv("SECURE_COOKIE_KEY"),
		},
	})
}

func sampleSitemapURLs() []lazyapp.SitemapURL {
	updated := time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC)
	return []lazyapp.SitemapURL{
		{Location: "/", LastUpdated: updated, ChangeFreq: "weekly", Priority: 1},
		{Location: "/posts", LastUpdated: updated, ChangeFreq: "weekly", Priority: 0.8},
		{Location: "/posts/hello-golazy", LastUpdated: updated, ChangeFreq: "monthly", Priority: 0.6},
		{Location: "/posts/embedded-content", LastUpdated: updated, ChangeFreq: "monthly", Priority: 0.6},
	}
}
