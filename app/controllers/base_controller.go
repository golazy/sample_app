package controllers

import (
	"context"
	"fmt"

	"golazy.dev/lazycontroller"
	"sample_app/app/services/timeservice"
)

type BaseController struct {
	lazycontroller.Base
	timeService timeservice.Service
}

func NewBaseController(ctx context.Context) (BaseController, error) {
	timeService, ok := timeservice.FromContext(ctx)
	if !ok {
		return BaseController{}, fmt.Errorf("time service is missing from application context")
	}
	base, err := lazycontroller.NewBase(ctx)
	if err != nil {
		return BaseController{}, err
	}

	controller := BaseController{
		Base:        base,
		timeService: timeService,
	}
	return controller, nil
}

func (c *BaseController) BeforeAction() error {
	c.Set("currentTime", c.timeService.Now().Format("2006-01-02 15:04:05 MST"))
	return nil
}
