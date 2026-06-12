package app

import (
	"embed"
	"io/fs"
)

// Files contains all templates and public assets required by the application.
//
//go:embed views public
var Files embed.FS

func Views() (fs.FS, error) {
	return fs.Sub(Files, "views")
}

func Public() (fs.FS, error) {
	return fs.Sub(Files, "public")
}
