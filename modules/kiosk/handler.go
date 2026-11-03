package kiosk

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/modules/frames"
	"mastergin/modules/sessions"
	"mastergin/modules/settings"
	"mastergin/modules/vouchers"
)

// Handler handles kiosk frontend and API.
type Handler struct {
	db  *gorm.DB
	cfg *core.Config
}

// NewHandler creates a new kiosk handler.
func NewHandler(db *gorm.DB, cfg *core.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// KioskPage renders the fullscreen interactive photobooth interface.
func (h *Handler) KioskPage(c *gin.Context) {
	cfgMap := settings.GetConfig(h.db, h.cfg.Port)
	activeFrames, _ := frames.ListActive(h.db)

	c.HTML(http.StatusOK, "kiosk", gin.H{
		"Config": cfgMap,
		"Frames": activeFrames,
	})
}

// FinishPayload is the JSON request payload sent by the kiosk canvas upon capture.
type FinishPayload struct {
	FrameID         uint     `json:"frame_id"`
	FilterApplied   string   `json:"filter_applied"`
	CompositeBase64 string   `json:"composite_base64"`
	RawPhotos       []string `json:"raw_photos"`
	PrintRequested  bool     `json:"print_requested"`
	VoucherCode     string   `json:"voucher_code,omitempty"`
}

// FinishSession handles the final composite upload, database storage, and print dispatch.
func (h *Handler) FinishSession(c *gin.Context) {
	var payload FinishPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	sessionCode := sessions.GenerateSessionCode()

	// 1. Save final composite image
	compositeFilename := fmt.Sprintf("%s_final.jpg", sessionCode)
	compositePathDisk := filepath.Join("storage", "photos", "final", compositeFilename)
	compositePathWeb := "/storage/photos/final/" + compositeFilename

	if err := sessions.SaveBase64Image(payload.CompositeBase64, compositePathDisk); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save composite photo: " + err.Error()})
		return
	}

	// 2. Create Session database record
	sess := sessions.Session{
		SessionCode:   sessionCode,
		FrameID:       payload.FrameID,
		CompositePath: compositePathWeb,
		FilterApplied: payload.FilterApplied,
		Status:        "completed",
		CreatedAt:     time.Now(),
	}
	if err := h.db.Create(&sess).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record session: " + err.Error()})
		return
	}

	// 3. Save raw photos
	for i, rawB64 := range payload.RawPhotos {
		rawFilename := fmt.Sprintf("%s_raw_%d.jpg", sessionCode, i+1)
		rawPathDisk := filepath.Join("storage", "photos", "raw", rawFilename)
		rawPathWeb := "/storage/photos/raw/" + rawFilename

		if err := sessions.SaveBase64Image(rawB64, rawPathDisk); err == nil {
			photoRecord := sessions.SessionPhoto{
				SessionID: sess.ID,
				ShotOrder: i + 1,
				PhotoPath: rawPathWeb,
				CreatedAt: time.Now(),
			}
			h.db.Create(&photoRecord)
		}
	}

	// 4. Handle Printing
	cfgMap := settings.GetConfig(h.db, h.cfg.Port)
	shouldPrint := (cfgMap.PrintEnabled && cfgMap.AutoPrint) || (cfgMap.PrintEnabled && payload.PrintRequested)
	if shouldPrint {
		_ = sessions.PrintPhoto(h.db, compositePathWeb, cfgMap.CopiesCount)
		sess.PrintedCount += cfgMap.CopiesCount
		h.db.Save(&sess)
	}

	// 5. Build QR Code Download URL
	// Uses the detected/configured LAN IP so mobile devices on the same WiFi can reach it
	qrDownloadURL := fmt.Sprintf("http://%s:%s/d/%s", cfgMap.LocalServerIP, h.cfg.Port, sessionCode)

	// 6. Redeem voucher if voucher was provided
	if strings.TrimSpace(payload.VoucherCode) != "" {
		_ = vouchers.ValidateAndRedeem(h.db, payload.VoucherCode, sessionCode)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"session_code":  sessionCode,
		"qr_url":        qrDownloadURL,
		"composite_url": compositePathWeb,
		"printed":       shouldPrint,
	})
}

// QRCode generates and returns a pure-Go QR Code PNG (100% offline).
func (h *Handler) QRCode(c *gin.Context) {
	text := c.Query("text")
	if strings.TrimSpace(text) == "" {
		c.String(http.StatusBadRequest, "Missing text parameter")
		return
	}

	pngBytes, err := qrcode.Encode(text, qrcode.Medium, 260)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate QR code: "+err.Error())
		return
	}

	c.Data(http.StatusOK, "image/png", pngBytes)
}
