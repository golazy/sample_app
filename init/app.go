package appinit

import (
	"golazy.dev/lazyapp"
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
		Helpers: []map[string]any{helpers.RegisterHelpers()},
	})
}
