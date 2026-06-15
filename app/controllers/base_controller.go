package controllers

import (
	"context"
	"fmt"

	"golazy.dev/lazycontroller"
	"sample_app/app/services/timeservice"
)

type BaseController struct {
	lazycontroller.Base
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
		Base: base,
	}
	controller.Set("currentTime", timeService.Now().Format("2006-01-02 15:04:05 MST"))
	return controller, nil
}
