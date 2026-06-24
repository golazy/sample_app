package appinit

import (
	"golazy.dev/lazyroutes"
	homecontroller "sample_app/app/controllers/home_controller"
)

func Draw(router *lazyroutes.Scope) {
	router.Resources(homecontroller.New, func(home *lazyroutes.Resource) {
		home.Singular("home")
		home.Plural("home")
		home.Path("")
	})
}
