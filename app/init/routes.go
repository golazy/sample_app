package appinit

import (
	"context"
	"net/http"

	lazycontroller "golazy.dev/controller"
	lazyroutes "golazy.dev/routes"
	"sample_app/app/controllers/home"
	postcontroller "sample_app/app/controllers/posts"
)

func Draw(ctx context.Context, mux *http.ServeMux) {
	mux.Handle("GET /{$}", lazycontroller.Bind(ctx, home.New, (*home.HomeController).Index))
	mux.Handle("GET /posts", lazycontroller.Bind(ctx, postcontroller.New, (*postcontroller.PostsController).Index))
	mux.Handle("GET /posts/{param}", lazycontroller.Bind(ctx, postcontroller.New, (*postcontroller.PostsController).Show))

	mux.Handle("/{$}", lazyroutes.MethodNotAllowed(http.MethodGet))
	mux.Handle("/posts", lazyroutes.MethodNotAllowed(http.MethodGet))
	mux.Handle("/posts/{param}", lazyroutes.MethodNotAllowed(http.MethodGet))
}
