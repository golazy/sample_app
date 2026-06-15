//go:build lazydev

package app

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestViewsUsesLocalTemplates(t *testing.T) {
	views, err := Views()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fs.Stat(views, "layouts/app.html.tpl"); err != nil {
		t.Fatalf("default layout missing from local views: %v", err)
	}
}

func TestViewsDirCanBeOverridden(t *testing.T) {
	dir := t.TempDir()
	if err := writeFile(filepath.Join(dir, "layouts", "app.html.tpl"), "layout"); err != nil {
		t.Fatal(err)
	}

	previous := ViewsDir
	ViewsDir = dir
	t.Cleanup(func() {
		ViewsDir = previous
	})

	views, err := Views()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fs.Stat(views, "layouts/app.html.tpl"); err != nil {
		t.Fatalf("default layout missing from overridden views dir: %v", err)
	}
}

func writeFile(filename string, content string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(content), 0o644)
}
