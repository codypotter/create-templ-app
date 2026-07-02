package http

import (
	"net/http"
	"sync/atomic"

	"github.com/a-h/templ"
	"github.com/codypotter/create-templ-app/internal/assets"
	"github.com/codypotter/create-templ-app/internal/views"
	"github.com/gin-gonic/gin"
)

var serverCount atomic.Int64

func Register(r *gin.Engine, resolver *assets.Resolver, distPath string) {
	asset := resolver.URL

	// Serves static assets locally. If ASSET_BASE_URL points elsewhere
	// (e.g. a CDN), these routes just go unused.
	static := r.Group("/static")
	static.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	})
	static.Static("/", distPath)

	r.GET("/", func(c *gin.Context) {
		render(c, http.StatusOK, views.Index(asset))
	})

	r.POST("/api/count", func(c *gin.Context) {
		render(c, http.StatusOK, views.Counter(serverCount.Add(1)))
	})
}

func render(c *gin.Context, status int, t templ.Component) {
	c.Status(status)
	if err := t.Render(c.Request.Context(), c.Writer); err != nil {
		_ = c.Error(err)
	}
}
