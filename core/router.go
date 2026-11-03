package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the Gin engine with global middleware,
// template loading, and static asset serving. Route groups for individual
// modules are registered by the caller (main.go) via their RegisterRoutes
// functions.
func SetupRouter(cfg *Config) *gin.Engine {
	if cfg.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware: Recovery (prevents crash on panic) and Logger.
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Load HTML templates using the custom multi-template renderer.
	r.HTMLRender = LoadTemplates("template/meridian")

	// Serve static assets (CSS, JS, images) from the meridian template folder.
	r.Static("/assets", "./template/meridian/assets")

	// Serve uploaded storage files (photos, frames).
	r.Static("/storage", "./storage")

	// Health-check endpoint.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}
