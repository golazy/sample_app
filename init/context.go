package appinit

import (
	"context"
	"fmt"

	postservice "sample_app/app/services/posts"
	"sample_app/app/services/timeservice"
)

// Context initializes the application dependencies and adds them to ctx.
// Embedded resource failures are programming errors, so startup fails fast.
func Context(ctx context.Context) context.Context {
	posts, err := postservice.New()
	if err != nil {
		panic(fmt.Errorf("initialize posts service: %w", err))
	}

	ctx = timeservice.WithContext(ctx, timeservice.New())
	ctx = postservice.WithContext(ctx, posts)
	return ctx
}
