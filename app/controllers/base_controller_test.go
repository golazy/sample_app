package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"golazy.dev/lazycontroller"
	"golazy.dev/lazyview"
	_ "golazy.dev/lazyview/gotmpl"
	"sample_app/app/services/timeservice"
)

type fixedTimeService struct {
	now time.Time
}

func (s fixedTimeService) Now() time.Time {
	return s.now
}

func TestBaseControllerSetsCurrentTime(t *testing.T) {
	views := fstest.MapFS{
		"layouts/app.html.tpl": {Data: []byte(`{{$content := .content}}{{$content}}`)},
		"home/index.html.tpl":  {Data: []byte(`{{.currentTime}}`)},
	}
	renderer, err := lazycontroller.NewRenderer(views)
	if err != nil {
		t.Fatal(err)
	}

	expected := time.Date(2026, time.June, 11, 12, 30, 0, 0, time.UTC)
	response := httptest.NewRecorder()
	ctx := context.Background()
	ctx = lazycontroller.WithRenderer(ctx, renderer)
	ctx = timeservice.WithContext(ctx, fixedTimeService{now: expected})

	controller, err := NewBaseController(ctx)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	if err := controller.BindRequest(response, request, lazyview.Route{Controller: "home"}); err != nil {
		t.Fatal(err)
	}
	if err := controller.BeforeAction(); err != nil {
		t.Fatal(err)
	}
	if err := controller.Render("index"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(response.Body.String(), "2026-06-11 12:30:00 UTC") {
		t.Fatalf("unexpected body: %q", response.Body.String())
	}
}
