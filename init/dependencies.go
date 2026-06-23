package appinit

import (
	"context"
	"fmt"

	"golazy.dev/lazydeps"
	postservice "sample_app/services/posts"
	"sample_app/services/timeservice"
)

// Dependencies initializes the application dependencies and adds them to ctx.
func Dependencies(deps *lazydeps.Scope) error {
	if _, err := lazydeps.Service(deps, "timeservice", func(ctx context.Context) (context.Context, timeservice.Service, error, context.CancelFunc) {
		service := timeservice.New()
		return timeservice.WithContext(ctx, service), service, nil, nil
	}); err != nil {
		return err
	}

	_, err := lazydeps.Service(deps, "posts", func(ctx context.Context) (context.Context, *postservice.Service, error, context.CancelFunc) {
		posts, err := postservice.New()
		if err != nil {
			return ctx, nil, fmt.Errorf("initialize posts service: %w", err), nil
		}

		return postservice.WithContext(ctx, posts), posts, nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}
