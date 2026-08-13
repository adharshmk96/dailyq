// Package web embeds the built DailyQ UI so a single binary serves both the
// SPA and the API.
//
// The contents of dist/ are produced by `nuxt generate` in ../dailyq-ui and
// copied here by `task ui` before `go build`. Only dist/.keep is committed, so
// the package still compiles when no UI build is present — the server then
// serves the API alone.
package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// FS returns the embedded build rooted at dist/.
func FS() (fs.FS, error) {
	return fs.Sub(dist, "dist")
}

// HasBuild reports whether a UI build was embedded.
func HasBuild() bool {
	_, err := dist.Open("dist/index.html")
	return err == nil
}
