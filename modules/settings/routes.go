package settings

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
)

// RegisterRoutes registers the settings routes.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *core.Config) {
	h := NewHandler(db, cfg)

	// Public endpoint for kiosk frontend to fetch settings
	r.GET("/api/settings", h.GetConfigJSON)

	// Protected admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth(cfg))
	{
		admin.GET("/settings", h.SettingsPage)
		admin.POST("/settings", h.UpdateSettings)
	}
}

