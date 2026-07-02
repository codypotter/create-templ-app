package http

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/codypotter/create-templ-app/internal/assets"
	"github.com/codypotter/create-templ-app/internal/views"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, resolver *assets.Resolver, distPath string) {
	asset := resolver.URL

	// Static assets — in prod, CloudFront/S3 handles this; locally Go serves it directly.
	static := r.Group("/static")
	static.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	})
	static.Static("/", distPath)

	r.GET("/", func(c *gin.Context) {
		render(c, http.StatusOK, views.Index(asset))
	})
}

func render(c *gin.Context, status int, t templ.Component) {
	c.Status(status)
	if err := t.Render(c.Request.Context(), c.Writer); err != nil {
		_ = c.Error(err)
	}
}
