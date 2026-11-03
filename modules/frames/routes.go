package frames

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
)

// RegisterRoutes registers the frame management routes.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *core.Config) {
	h := NewHandler(db, cfg)

	// Public API for kiosk & slot editor
	r.GET("/api/frames", h.APIListActive)
	r.GET("/api/frames/:id", h.APIGetFrame)

	// Admin protected routes
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth(cfg))
	{
		admin.GET("/frames", h.AdminFramesPage)
		admin.POST("/frames", h.CreateFrame)
		admin.POST("/frames/:id/update", h.UpdateFrame)
		admin.POST("/frames/:id/toggle", h.ToggleFrame)
		admin.POST("/frames/:id/delete", h.DeleteFrame)
	}
}

