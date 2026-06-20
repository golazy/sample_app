package controllers

import (
	"context"
	"fmt"

	"golazy.dev/lazycontroller"
	"golazy.dev/lazysession"
	"sample_app/app/services/timeservice"
)

const visitCounterKey = "visit_count"

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
	session, ok, err := c.session()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	if err := c.trackVisits(session); err != nil {
		return err
	}
	c.Set("flashMessages", session.Flashes())
	return nil
}

func (c *BaseController) session() (*lazysession.Session, bool, error) {
	request := c.Request()
	if request == nil {
		return nil, false, nil
	}
	if _, ok := lazysession.ManagerFromContext(request.Context()); !ok {
		return nil, false, nil
	}
	session, err := lazysession.Get(request)
	if err != nil {
		return nil, true, err
	}
	return session, true, nil
}

func (c *BaseController) Session() (*lazysession.Session, bool, error) {
	return c.session()
}

func (c *BaseController) trackVisits(session *lazysession.Session) error {
	count, err := visitCounterValue(session.Values[visitCounterKey])
	if err != nil {
		return err
	}
	count++
	session.Values[visitCounterKey] = count
	c.Set("visitCount", count)
	return nil
}

func visitCounterValue(value any) (int, error) {
	if value == nil {
		return 0, nil
	}

	switch count := value.(type) {
	case int:
		return count, nil
	case int8:
		return int(count), nil
	case int16:
		return int(count), nil
	case int32:
		return int(count), nil
	case int64:
		return int(count), nil
	case uint:
		return int(count), nil
	case uint8:
		return int(count), nil
	case uint16:
		return int(count), nil
	case uint32:
		return int(count), nil
	case uint64:
		return int(count), nil
	case float32:
		return int(count), nil
	case float64:
		return int(count), nil
	default:
		return 0, fmt.Errorf("visit counter session value has type %T, want integer", value)
	}
}
