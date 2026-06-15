//go:build !lazydev

package app

import (
	"io/fs"
	"testing"
)

func TestViewsUsesEmbeddedTemplates(t *testing.T) {
	views, err := Views()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fs.Stat(views, "layouts/app.html.tpl"); err != nil {
		t.Fatalf("default layout missing from embedded views: %v", err)
	}
}
