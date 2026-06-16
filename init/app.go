package appinit

import (
	"golazy.dev/lazyapp"
	"golazy.dev/lazysession"
	_ "golazy.dev/lazyview/gotmpl"
	"sample_app/app"
	"sample_app/app/helpers"
)

// In production, keep this secret stable and load it from
// os.Getenv("SECURE_COOKIE_KEY") or another secret manager.
const secureCookieKey = "a5c9d0bd530c48964911825a1914bf9a892e038c3474914b2792c1a5def3aa15"

func App() *lazyapp.App {
	return lazyapp.New(lazyapp.Config{
		Name:    "sample_app",
		Drawer:  Draw,
		Public:  app.Public,
		Views:   app.Views,
		Context: Context,
		Helpers: []map[string]any{helpers.RegisterHelpers()},
		Sessions: lazysession.Config{
			Name: "sample_app_session",
			KeyPairs: [][]byte{
				[]byte(secureCookieKey),
			},
		},
	})
}
