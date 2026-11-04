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
	r.GET("/live-wall", h.LiveWallPage)
	r.GET("/api/live-wall/recent", h.LiveWallRecentAPI)
	r.GET("/api/qr", h.QRCode)
	r.POST("/api/kiosk/finish", h.FinishSession)

	// Self-service vending payment endpoints
	r.POST("/api/kiosk/payment/create", h.CreatePaymentOrder)
	r.GET("/api/kiosk/payment/status", h.CheckPaymentStatus)
	r.POST("/api/kiosk/payment/simulate", h.SimulatePaymentSuccess)

	// Kiosk heartbeat
	r.POST("/api/kiosk/heartbeat", h.KioskHeartbeatAPI)

	// Remote operator mobile controller
	r.GET("/operator", h.OperatorPage)
	r.GET("/api/operator/status", h.OperatorStatusAPI)
	r.POST("/api/operator/command", h.OperatorCommandAPI)
	r.POST("/api/operator/reprint-last", h.OperatorReprintLastAPI)
	r.POST("/api/operator/reset-paper", h.OperatorResetPaperAPI)
	r.POST("/api/operator/quick-voucher", h.OperatorQuickVoucherAPI)
}
