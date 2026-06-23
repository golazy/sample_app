package appinit

import (
	"context"
	"fmt"

	postservice "sample_app/services/posts"
	"sample_app/services/timeservice"
)

// Context initializes the application dependencies and adds them to ctx.
func Context(ctx context.Context) (context.Context, error) {
	posts, err := postservice.New()
	if err != nil {
		return ctx, fmt.Errorf("initialize posts service: %w", err)
	}

	ctx = timeservice.WithContext(ctx, timeservice.New())
	ctx = postservice.WithContext(ctx, posts)
	return ctx, nil
}
