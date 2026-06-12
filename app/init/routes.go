package appinit

import (
	"context"
	"net/http"

	"golazy.dev/lazyroutes"
	"sample_app/app/controllers/home"
	postcontroller "sample_app/app/controllers/posts"
)

func Draw(ctx context.Context, mux *http.ServeMux) {
	mux.Handle("GET /{$}", lazyroutes.Bind(ctx, home.New, (*home.HomeController).Index))
	lazyroutes.Resources(ctx, mux, postcontroller.New)
}
