package helloworldservice

import "context"

type Service interface {
	Hello() string
}

type service struct{}

type contextKey struct{}

func New(ctx context.Context) (context.Context, Service, error, context.CancelFunc) {
	service := service{}
	return WithContext(ctx, service), service, nil, nil
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
