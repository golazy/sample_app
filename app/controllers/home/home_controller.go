package home

import (
	"context"
	"net/http"

	"sample_app/app/controllers"
)

type HomeController struct {
	controllers.BaseController
}

func New(ctx context.Context) (*HomeController, error) {
	base, err := controllers.NewBaseController(ctx)
	if err != nil {
		return nil, err
	}
	return &HomeController{BaseController: base}, nil
}

func (c *HomeController) Index(_ http.ResponseWriter, _ *http.Request) error {
	c.Set("title", "Home")
	return c.Render("index")
}
