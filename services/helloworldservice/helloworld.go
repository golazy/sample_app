package helloworldservice

import "context"

type Service interface {
	Hello() string
}

type service struct{}

type contextKey struct{}

func New() Service {
	return service{}
}

func (service) Hello() string {
	return "Hello"
}

func WithContext(ctx context.Context, service Service) context.Context {
	return context.WithValue(ctx, contextKey{}, service)
}

func FromContext(ctx context.Context) (Service, bool) {
	service, ok := ctx.Value(contextKey{}).(Service)
	return service, ok
}
