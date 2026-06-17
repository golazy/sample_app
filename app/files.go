package app

import (
	"embed"

	"golazy.dev/lazyapp"
)

// Files contains all templates and public assets required by the application.
//
//go:embed views public
var Files embed.FS

var Views = lazyapp.MustSub(Files, "views")
var Public = lazyapp.MustSub(Files, "public")
