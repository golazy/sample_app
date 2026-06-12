package appinit

import (
	"context"
	"fmt"
	"net/http"

	"golazy.dev/lazycontroller"
	"golazy.dev/lazyroutes"
	"sample_app/app"
	postservice "sample_app/app/services/posts"
	"sample_app/app/services/timeservice"
)

// Context initializes the application dependencies and adds them to ctx.
// Embedded resource failures are programming errors, so startup fails fast.
func Context(ctx context.Context) context.Context {
	views, err := app.Views()
	if err != nil {
		panic(fmt.Errorf("open embedded views: %w", err))
	}
	renderer, err := lazycontroller.NewRenderer(views)
	if err != nil {
		panic(fmt.Errorf("initialize renderer: %w", err))
	}

	posts, err := postservice.New()
	if err != nil {
		panic(fmt.Errorf("initialize posts service: %w", err))
	}

	public, err := app.Public()
	if err != nil {
		panic(fmt.Errorf("open embedded public files: %w", err))
	}

	ctx = lazycontroller.WithRenderer(ctx, renderer)
	ctx = timeservice.WithContext(ctx, timeservice.New())
	ctx = postservice.WithContext(ctx, posts)
	ctx = lazyroutes.WithPublic(ctx, http.FileServerFS(public))
	return ctx
}
