package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"prismeda/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(svc *service.Service, webDir string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/api/health", Health(svc))
	r.POST("/api/datasets", UploadDataset(svc))
	r.POST("/api/datasets/sample", LoadSample(svc))
	r.GET("/api/datasets/:id", GetProfile(svc))
	r.GET("/api/datasets/:id/preview", GetPreview(svc))
	r.GET("/api/datasets/:id/scatter", GetScatter(svc))

	index := filepath.Join(webDir, "index.html")
	if _, err := os.Stat(index); err != nil {
		r.NoRoute(func(c *gin.Context) { fail(c, http.StatusNotFound, "frontend not built — run 'npm run build' in web/") })
		return r
	}

	r.Static("/assets", filepath.Join(webDir, "assets"))
	if _, err := os.Stat(filepath.Join(webDir, "favicon.svg")); err == nil {
		r.StaticFile("/favicon.svg", filepath.Join(webDir, "favicon.svg"))
	}
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			fail(c, http.StatusNotFound, "endpoint not found")
			return
		}
		c.File(index) // SPA history fallback
	})
	return r
}
