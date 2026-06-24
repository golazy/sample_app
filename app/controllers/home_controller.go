package controllers

import (
	"context"
	"fmt"
	"net/http"

	"sample_app/services/helloworldservice"

	"golazy.dev/lazycontroller"
)

type HomeController struct {
	lazycontroller.Base
	helloService helloworldservice.Service
}

func NewHomeController(ctx context.Context) (*HomeController, error) {
	base, err := lazycontroller.NewBase(ctx)
	if err != nil {
		return nil, err
	}
	helloService, ok := helloworldservice.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("hello world service is missing from application context")
	}
	return &HomeController{Base: base, helloService: helloService}, nil
}

func (c *HomeController) Index(_ http.ResponseWriter, _ *http.Request) error {
	c.Set("title", c.helloService.Hello())
	return nil
}
