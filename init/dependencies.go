package appinit

import (
	"context"

	"golazy.dev/lazydeps"
	"sample_app/services/helloworldservice"
)

// Dependencies initializes the application dependencies and adds them to ctx.
func Dependencies(deps *lazydeps.Scope) error {
	_, err := lazydeps.Service(deps, "helloworldservice", func(ctx context.Context) (context.Context, helloworldservice.Service, error, context.CancelFunc) {
		service := helloworldservice.New()
		return helloworldservice.WithContext(ctx, service), service, nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}
