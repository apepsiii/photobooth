package vouchers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
)

// RegisterRoutes registers voucher routes.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *core.Config) {
	h := NewHandler(db, cfg)

	// Public API for kiosk validation
	r.POST("/api/vouchers/validate", h.APIValidate)

	// Protected admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth(cfg))
	{
		admin.GET("/vouchers", h.AdminVouchersPage)
		admin.POST("/vouchers/generate", h.Generate)
		admin.POST("/vouchers/toggle-mode", h.ToggleMode)
		admin.POST("/vouchers/:id/delete", h.Delete)
	}
}

