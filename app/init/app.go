package appinit

import (
	"golazy.dev/lazyapp"
	"sample_app/app"
)

func App() *lazyapp.App {
	return lazyapp.New(lazyapp.Config{
		Name:    "sample_app",
		Drawer:  Draw,
		Public:  app.Public,
		Views:   app.Views,
		Context: Context,
	})
}
