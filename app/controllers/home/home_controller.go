package home

import (
	"context"
	"net/http"

	"sample_app/app/controllers"
)

type HomeController struct {
	controllers.BaseController
}

type homeMetadata struct{}

func (homeMetadata) Title() string {
	return "Home"
}

func (homeMetadata) Canonical() string {
	return "/"
}

func New(ctx context.Context) (*HomeController, error) {
	base, err := controllers.NewBaseController(ctx)
	if err != nil {
		return nil, err
	}
	return &HomeController{BaseController: base}, nil
}

func (c *HomeController) Index(_ http.ResponseWriter, _ *http.Request) error {
	c.Metadata(homeMetadata{})
	return nil
}

func (c *HomeController) Flash(_ http.ResponseWriter, _ *http.Request) error {
	session, ok, err := c.Session()
	if err != nil {
		return err
	}
	if ok {
		session.AddFlash("This is a sample flash message")
	}
	return c.RedirectTo("/")
}
