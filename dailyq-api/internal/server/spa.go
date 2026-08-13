package server

import (
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"dailyq-api/internal/web"
)

// spaHandler serves the embedded UI build: real files are served as-is, and
// every other path falls back to index.html so client-side routing works on a
// hard refresh or a deep link.
func spaHandler(log *slog.Logger) (gin.HandlerFunc, error) {
	assets, err := web.FS()
	if err != nil {
		return nil, err
	}

	if !web.HasBuild() {
		log.Warn("serving placeholder UI: no ui build embedded, run `task build`")
	}

	fileServer := http.FileServer(http.FS(assets))

	return func(c *gin.Context) {
		req := c.Request

		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			c.JSON(http.StatusNotFound, notFoundBody())
			return
		}

		name := strings.TrimPrefix(path.Clean(req.URL.Path), "/")
		if name == "" || name == "." {
			serveIndex(c, assets)
			return
		}

		f, err := assets.Open(name)
		if err != nil {
			// Unknown path: let the SPA router decide what to render.
			serveIndex(c, assets)
			return
		}
		info, statErr := f.Stat()
		_ = f.Close()
		if statErr != nil || info.IsDir() {
			serveIndex(c, assets)
			return
		}

		// Nuxt emits content-hashed filenames under /_nuxt, so they are safe to
		// cache forever.
		if strings.HasPrefix(name, "_nuxt/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}

		fileServer.ServeHTTP(c.Writer, req)
	}, nil
}

func serveIndex(c *gin.Context, assets fs.FS) {
	body, err := fs.ReadFile(assets, "index.html")
	if err != nil {
		c.JSON(http.StatusNotFound, notFoundBody())
		return
	}

	// index.html must never be cached, or clients pin themselves to a stale
	// bundle after a deploy.
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", body)
}

func notFoundBody() gin.H {
	return gin.H{"error": gin.H{
		"code":    "not_found",
		"message": "resource not found",
	}}
}
