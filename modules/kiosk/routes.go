package kiosk

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
)

// RegisterRoutes registers kiosk photobooth frontend and API routes.
func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *core.Config) {
	h := NewHandler(db, cfg)

	r.GET("/kiosk", h.KioskPage)
	r.GET("/api/qr", h.QRCode)
	r.POST("/api/kiosk/finish", h.FinishSession)
}
