package timeservice

import (
	"context"
	"time"
)

type Service interface {
	Now() time.Time
}

type service struct{}

func New() Service {
	return service{}
}

func (service) Now() time.Time {
	return time.Now()
}

type contextKey struct{}

func WithContext(ctx context.Context, service Service) context.Context {
	return context.WithValue(ctx, contextKey{}, service)
}

func FromContext(ctx context.Context) (Service, bool) {
	service, ok := ctx.Value(contextKey{}).(Service)
	return service, ok
}
