package appinit

import (
	"golazy.dev/lazydeps"
	"sample_app/services/helloworldservice"
)

// Dependencies initializes the application dependencies and adds them to ctx.
func Dependencies(deps *lazydeps.Scope) error {
	_, err := lazydeps.Service(deps, "helloworldservice", helloworldservice.New)
	return err
}
