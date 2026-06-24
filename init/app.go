package appinit

import (
	"time"

	"golazy.dev/lazyapp"
	"golazy.dev/lazycontrolplane"
	"golazy.dev/lazyseo"
	_ "golazy.dev/lazyview/gotmpl"
	"sample_app/app"
)

func App() *lazyapp.App {
	return lazyapp.New(lazyapp.Config{
		Name:         "sample_app",
		Drawer:       Draw,
		Public:       app.Public,
		Views:        app.Views,
		Dependencies: Dependencies,
		ControlPlane: lazycontrolplane.Config{},
		SEO: []lazyseo.Option{
			lazyseo.SiteName("GoLazy"),
			lazyseo.Description("A small GoLazy sample application."),
			lazyseo.Language("en"),
			lazyseo.Type("website"),
			lazyseo.Locale("en_US"),
			lazyseo.TwitterCardType("summary"),
		},
		Sitemap: lazyapp.SitemapConfig{
			URLs: sampleSitemapURLs(),
		},
	})
}

func sampleSitemapURLs() []lazyapp.SitemapURL {
	updated := time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC)
	return []lazyapp.SitemapURL{
		{Location: "/", LastUpdated: updated, ChangeFreq: "weekly", Priority: 1},
	}
}
