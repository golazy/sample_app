package homecontroller

import (
	"context"
	"fmt"
	"net/http"

	"sample_app/app/controllers"
	"sample_app/services/helloworldservice"
)

type HomeController struct {
	controllers.BaseController
	helloService helloworldservice.Service
}

func New(ctx context.Context) (*HomeController, error) {
	base, err := controllers.NewBaseController(ctx)
	if err != nil {
		return nil, err
	}
	helloService, ok := helloworldservice.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("hello world service is missing from application context")
	}
	return &HomeController{BaseController: base, helloService: helloService}, nil
}

func (c *HomeController) Index(_ http.ResponseWriter, _ *http.Request) error {
	c.Set("title", c.helloService.Hello())
	return nil
}
