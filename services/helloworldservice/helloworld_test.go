package helloworldservice

import (
	"context"
	"testing"
)

func TestServiceSaysHello(t *testing.T) {
	service := New()
	if got, want := service.Hello(), "Hello"; got != want {
		t.Fatalf("Hello() = %q, want %q", got, want)
	}
}

func TestContextHelpers(t *testing.T) {
	service := New()
	ctx := WithContext(context.Background(), service)

	got, ok := FromContext(ctx)
	if !ok {
		t.Fatal("service missing from context")
	}
	if got.Hello() != "Hello" {
		t.Fatalf("context service Hello() = %q", got.Hello())
	}
}
