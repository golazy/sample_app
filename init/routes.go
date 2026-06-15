package appinit

import (
	"golazy.dev/lazyroutes"
	"sample_app/app/controllers/home"
	postcontroller "sample_app/app/controllers/posts"
)

func Draw(router *lazyroutes.Scope) {
	router.Get("/", home.New, (*home.HomeController).Index)
	router.Resources(postcontroller.New)
}
