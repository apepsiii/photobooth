package sessions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
)

// RegisterRoutes registers the session management and public download routes.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *core.Config) {
	h := NewHandler(db, cfg)

	// Public routes for smartphone QR softfile download
	r.GET("/d/:code", h.DownloadSoftfilePage)
	r.GET("/d/:code/download", h.DownloadPhoto)
	r.GET("/d/:code/gif", h.DownloadGIF)

	// Protected admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth(cfg))
	{
		admin.GET("/sessions", h.AdminSessionsPage)
		admin.GET("/sessions/export-zip", h.ExportAllZip)
		admin.GET("/sessions/export-csv", h.ExportCSV)
		admin.POST("/sessions/:id/reprint", h.ReprintSession)
		admin.POST("/sessions/:id/delete", h.DeleteSession)
		admin.POST("/sessions/test-print", h.TestPrint)
	}
}

