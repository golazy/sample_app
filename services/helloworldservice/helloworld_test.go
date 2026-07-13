package helloworldservice

import (
	"context"
	"testing"
)

func TestServiceSaysHello(t *testing.T) {
	_, service, err, stop := New(context.Background())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if stop != nil {
		defer stop()
	}
	if got, want := service.Hello(), "Hello"; got != want {
		t.Fatalf("Hello() = %q, want %q", got, want)
	}
}

func TestContextHelpers(t *testing.T) {
	ctx, _, err, stop := New(context.Background())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if stop != nil {
		defer stop()
	}

	got, ok := FromContext(ctx)
	if !ok {
		t.Fatal("service missing from context")
	}
	if got.Hello() != "Hello" {
		t.Fatalf("context service Hello() = %q", got.Hello())
	}
}
