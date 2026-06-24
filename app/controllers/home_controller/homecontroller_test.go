package homecontroller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"golazy.dev/lazycontroller"
	"golazy.dev/lazyview"
	_ "golazy.dev/lazyview/gotmpl"
	"sample_app/services/helloworldservice"
)

func TestHomeControllerIndexSetsTitleFromHelloService(t *testing.T) {
	views := fstest.MapFS{
		"layouts/app.html.tpl": {Data: []byte(`{{.content}}`)},
		"home/index.html.tpl":  {Data: []byte(`title={{.title}}`)},
	}
	renderer, err := lazycontroller.NewRenderer(views)
	if err != nil {
		t.Fatal(err)
	}

	ctx := lazycontroller.WithRenderer(context.Background(), renderer)
	ctx = helloworldservice.WithContext(ctx, stubHelloService{message: "Hello from test"})

	controller, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	if err := controller.BindRequest(response, request, lazyview.Route{Controller: "home", Action: "Index"}); err != nil {
		t.Fatal(err)
	}

	if err := controller.Index(response, request); err != nil {
		t.Fatal(err)
	}
	if err := controller.Render("index"); err != nil {
		t.Fatal(err)
	}
	if got, want := response.Body.String(), "title=Hello from test"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHomeControllerRequiresHelloService(t *testing.T) {
	renderer, err := lazycontroller.NewRenderer(fstest.MapFS{
		"layouts/app.html.tpl": {Data: []byte(`{{.content}}`)},
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := lazycontroller.WithRenderer(context.Background(), renderer)
	_, err = New(ctx)
	if err == nil {
		t.Fatal("expected missing service error")
	}
	if !strings.Contains(err.Error(), "hello world service is missing") {
		t.Fatalf("error = %v, want missing hello service", err)
	}
}

type stubHelloService struct {
	message string
}

func (s stubHelloService) Hello() string {
	return s.message
}
