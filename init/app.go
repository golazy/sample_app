package appinit

import (
	"os"

	"golazy.dev/lazyapp"
	"golazy.dev/lazysession"
	_ "golazy.dev/lazyview/gotmpl"
	"sample_app/app"
	"sample_app/app/helpers"
)

func App() *lazyapp.App {
	return lazyapp.New(lazyapp.Config{
		Name:    "sample_app",
		Drawer:  Draw,
		Public:  app.Public,
		Views:   app.Views,
		Context: Context,
		Helpers: lazyapp.Helpers{helpers.RegisterHelpers()},
		Sessions: lazysession.Config{
			Key: os.Getenv("SECURE_COOKIE_KEY"),
		},
	})
}
