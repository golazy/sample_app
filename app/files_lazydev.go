//go:build lazydev

package app

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ViewsDir is the local template directory used by lazydev builds.
var ViewsDir = "app/views"

// Files contains embedded public assets for lazydev builds.
//
//go:embed public
var Files embed.FS

func Views() (fs.FS, error) {
	if ViewsDir == "" {
		return nil, fmt.Errorf("views directory is required")
	}
	dir, err := resolveViewsDir(ViewsDir)
	if err != nil {
		return nil, err
	}
	return os.DirFS(dir), nil
}

func Public() (fs.FS, error) {
	return fs.Sub(Files, "public")
}

func resolveViewsDir(dir string) (string, error) {
	if filepath.IsAbs(dir) {
		if viewsDirExists(dir) {
			return dir, nil
		}
		return "", fmt.Errorf("views directory %q does not contain layouts/app.html.tpl", dir)
	}

	if viewsDirExists(dir) {
		return dir, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	for {
		candidate := filepath.Join(cwd, dir)
		if viewsDirExists(candidate) {
			return candidate, nil
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}

	return "", fmt.Errorf("views directory %q does not contain layouts/app.html.tpl", dir)
}

func viewsDirExists(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "layouts", "app.html.tpl"))
	return err == nil && !info.IsDir()
}
