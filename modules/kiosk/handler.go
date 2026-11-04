package kiosk

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
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
	db         *gorm.DB
	cfg        *core.Config
	stateMu    sync.Mutex
	lastSeen   time.Time
	currScreen string
	pendingCmd string
}

// NewHandler creates a new kiosk handler.
func NewHandler(db *gorm.DB, cfg *core.Config) *Handler {
	return &Handler{
		db:         db,
		cfg:        cfg,
		currScreen: "screenWelcome",
	}
}

// KioskPage renders the fullscreen interactive photobooth interface.
func (h *Handler) KioskPage(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

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
	LivePhotoBase64 string   `json:"live_photo_base64,omitempty"`
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

	// 1b. Save Live Photo (MP4) if provided
	livePhotoPathWeb := ""
	if strings.TrimSpace(payload.LivePhotoBase64) != "" {
		liveFilename := fmt.Sprintf("%s_live.mp4", sessionCode)
		livePathDisk := filepath.Join("storage", "photos", "final", liveFilename)
		if err := sessions.SaveBase64File(payload.LivePhotoBase64, livePathDisk); err == nil {
			livePhotoPathWeb = "/storage/photos/final/" + liveFilename
		}
	}

	// Calculate daily queue number (antrean)
	var todayCount int64
	todayStart := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&sessions.Session{}).Where("created_at >= ?", todayStart).Count(&todayCount)
	queueNum := int(todayCount) + 1

	// 2. Create Session database record
	sess := sessions.Session{
		SessionCode:   sessionCode,
		QueueNumber:   queueNum,
		FrameID:       payload.FrameID,
		CompositePath: compositePathWeb,
		LivePhotoPath: livePhotoPathWeb,
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
		"success":        true,
		"session_code":   sessionCode,
		"queue_number":   queueNum,
		"qr_url":         qrDownloadURL,
		"composite_url":  compositePathWeb,
		"live_photo_url": livePhotoPathWeb,
		"printed":        shouldPrint,
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

// LiveWallPage renders the live projector and digital photowall interface.
func (h *Handler) LiveWallPage(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	cfgMap := settings.GetConfig(h.db, h.cfg.Port)

	var todayCount int64
	todayStart := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&sessions.Session{}).Where("status = ? AND created_at >= ?", "completed", todayStart).Count(&todayCount)

	c.HTML(http.StatusOK, "live-wall", gin.H{
		"Config":     cfgMap,
		"TodayCount": todayCount,
	})
}

// LiveSessionItem is the JSON item returned to the live wall.
type LiveSessionItem struct {
	ID            uint   `json:"id"`
	SessionCode   string `json:"session_code"`
	QueueNumber   int    `json:"queue_number"`
	CompositePath string `json:"composite_path"`
	LivePhotoPath string `json:"live_photo_path,omitempty"`
	FilterApplied string `json:"filter_applied"`
	TimeAgo       string `json:"time_ago"`
	CreatedAt     string `json:"created_at"`
}

// LiveWallRecentAPI returns recent completed photo sessions for live projector displays.
func (h *Handler) LiveWallRecentAPI(c *gin.Context) {
	cfgMap := settings.GetConfig(h.db, h.cfg.Port)

	var sessList []sessions.Session
	h.db.Where("status = ?", "completed").
		Order("id desc").
		Limit(24).
		Find(&sessList)

	var todayCount int64
	todayStart := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&sessions.Session{}).Where("status = ? AND created_at >= ?", "completed", todayStart).Count(&todayCount)

	var totalCount int64
	h.db.Model(&sessions.Session{}).Where("status = ?", "completed").Count(&totalCount)

	items := make([]LiveSessionItem, 0, len(sessList))
	now := time.Now()
	for _, s := range sessList {
		diff := now.Sub(s.CreatedAt)
		timeAgo := "Baru saja"
		if diff.Minutes() >= 1 && diff.Minutes() < 60 {
			timeAgo = fmt.Sprintf("%d menit lalu", int(diff.Minutes()))
		} else if diff.Hours() >= 1 && diff.Hours() < 24 {
			timeAgo = fmt.Sprintf("%d jam lalu", int(diff.Hours()))
		} else if diff.Hours() >= 24 {
			timeAgo = fmt.Sprintf("%d hari lalu", int(diff.Hours()/24))
		}

		items = append(items, LiveSessionItem{
			ID:            s.ID,
			SessionCode:   s.SessionCode,
			QueueNumber:   s.QueueNumber,
			CompositePath: s.CompositePath,
			LivePhotoPath: s.LivePhotoPath,
			FilterApplied: s.FilterApplied,
			TimeAgo:       timeAgo,
			CreatedAt:     s.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"event_name":  cfgMap.EventName,
		"school_name": cfgMap.SchoolName,
		"today_count": todayCount,
		"total_count": totalCount,
		"sessions":    items,
	})
}

// CreatePaymentPayload is the request payload to initiate self-service payment.
type CreatePaymentPayload struct {
	PackageType string `json:"package_type"` // "single", "double"
}

// CreatePaymentOrder initiates a self-service vending payment transaction.
func (h *Handler) CreatePaymentOrder(c *gin.Context) {
	var payload CreatePaymentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		payload.PackageType = "single"
	}

	cfgMap := settings.GetConfig(h.db, h.cfg.Port)

	amount := cfgMap.PriceSingle
	copies := 1
	if payload.PackageType == "double" {
		amount = cfgMap.PriceDouble
		copies = 2
	}

	orderID := fmt.Sprintf("ORD-%d-%04d", time.Now().Unix(), time.Now().Nanosecond()%10000)

	order := PaymentOrder{
		OrderID:     orderID,
		PackageType: payload.PackageType,
		Amount:      amount,
		Copies:      copies,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if err := h.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment order: " + err.Error()})
		return
	}

	// Build QRIS payload string: use merchant static code or generate dynamic payment string
	qrisData := cfgMap.QRISStaticCode
	if strings.TrimSpace(qrisData) == "" {
		qrisData = fmt.Sprintf("00020101021226580014ID.GO.BOOTH.WWW011893600912%s520458125303360540%d5802ID5913%s6007JAKARTA6304", orderID, amount, cfgMap.AppName)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"order_id":     orderID,
		"package_type": payload.PackageType,
		"amount":       amount,
		"copies":       copies,
		"status":       "pending",
		"qris_data":    qrisData,
		"qr_url":       "/api/qr?text=" + strings.ReplaceAll(qrisData, " ", "%20"),
	})
}

// CheckPaymentStatus polls the payment status of an order.
func (h *Handler) CheckPaymentStatus(c *gin.Context) {
	orderID := c.Query("order_id")
	if strings.TrimSpace(orderID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order_id"})
		return
	}

	var order PaymentOrder
	if err := h.db.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"order_id":     order.OrderID,
		"status":       order.Status,
		"package_type": order.PackageType,
		"copies":       order.Copies,
		"amount":       order.Amount,
	})
}

// SimulatePaymentSuccess approves a pending order for testing or operator cash reception.
func (h *Handler) SimulatePaymentSuccess(c *gin.Context) {
	var body struct {
		OrderID string `json:"order_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.OrderID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order_id"})
		return
	}

	var order PaymentOrder
	if err := h.db.Where("order_id = ?", body.OrderID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	now := time.Now()
	order.Status = "paid"
	order.PaidAt = &now
	h.db.Save(&order)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment confirmed successfully",
		"order":   order,
	})
}

// KioskHeartbeatAPI receives heartbeat from active kiosk display and delivers pending commands.
func (h *Handler) KioskHeartbeatAPI(c *gin.Context) {
	var body struct {
		Screen string `json:"screen"`
	}
	_ = c.ShouldBindJSON(&body)

	h.stateMu.Lock()
	h.lastSeen = time.Now()
	if body.Screen != "" {
		h.currScreen = body.Screen
	}
	cmd := h.pendingCmd
	h.pendingCmd = ""
	h.stateMu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"command": cmd,
	})
}

// OperatorPage renders the mobile-first touch control panel for on-ground booth staff.
func (h *Handler) OperatorPage(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	cfgMap := settings.GetConfig(h.db, h.cfg.Port)

	c.HTML(http.StatusOK, "operator", gin.H{
		"Config": cfgMap,
	})
}

// OperatorStatusAPI returns live state of kiosk, paper roll gauge, and last session snapshot.
func (h *Handler) OperatorStatusAPI(c *gin.Context) {
	cfgMap := settings.GetConfig(h.db, h.cfg.Port)

	h.stateMu.Lock()
	isOnline := !h.lastSeen.IsZero() && time.Since(h.lastSeen) < 6*time.Second
	lastSeenAgo := "Belum aktif"
	if !h.lastSeen.IsZero() {
		lastSeenAgo = fmt.Sprintf("%d dtk lalu", int(time.Since(h.lastSeen).Seconds()))
	}
	currScreen := h.currScreen
	h.stateMu.Unlock()

	// Screen human labels
	screenLabels := map[string]string{
		"screenWelcome": "Standby (Layar Awal)",
		"screenVoucher": "Input Tiket Voucher",
		"screenPayment": "Menunggu Bayar QRIS",
		"screenFrames":  "Memilih Frame",
		"screenCapture": "Sesi Kamera Berlangsung",
		"screenReview":  "Editor Doodle & Stiker",
		"screenResult":  "Selesai & Cetak Foto",
	}
	screenText := currScreen
	if lbl, ok := screenLabels[currScreen]; ok {
		screenText = lbl
	}

	// Fetch most recent session
	var lastSess sessions.Session
	var lastSessData gin.H
	if err := h.db.Order("id desc").First(&lastSess).Error; err == nil {
		lastSessData = gin.H{
			"code":          lastSess.SessionCode,
			"image_url":     lastSess.CompositePath,
			"created_at":    lastSess.CreatedAt.Format("15:04:05"),
			"printed_count": lastSess.PrintedCount,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"online":          isOnline,
		"last_seen_ago":   lastSeenAgo,
		"current_screen":  currScreen,
		"screen_label":    screenText,
		"paper_capacity":  cfgMap.PaperCapacity,
		"paper_used":      cfgMap.PaperUsed,
		"paper_remaining": cfgMap.PaperRemaining,
		"last_session":    lastSessData,
	})
}

// OperatorCommandAPI queues an urgent action (shutter trigger, reset to standby) for the kiosk.
func (h *Handler) OperatorCommandAPI(c *gin.Context) {
	var body struct {
		Command string `json:"command"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Command) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Perintah tidak valid"})
		return
	}

	h.stateMu.Lock()
	h.pendingCmd = body.Command
	h.stateMu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"command": body.Command,
	})
}

// OperatorReprintLastAPI triggers an immediate reprint of the most recent session composite.
func (h *Handler) OperatorReprintLastAPI(c *gin.Context) {
	var lastSess sessions.Session
	if err := h.db.Order("id desc").First(&lastSess).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Belum ada sesi foto yang tersimpan"})
		return
	}

	cfgMap := settings.GetConfig(h.db, h.cfg.Port)
	copies := cfgMap.CopiesCount
	if copies <= 0 {
		copies = 1
	}

	if err := sessions.PrintPhoto(h.db, lastSess.CompositePath, copies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim cetak ke printer: " + err.Error()})
		return
	}

	lastSess.PrintedCount += copies
	h.db.Save(&lastSess)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      fmt.Sprintf("Berhasil mencetak ulang sesi %s (%d lembar)", lastSess.SessionCode, copies),
		"session_code": lastSess.SessionCode,
	})
}

// OperatorResetPaperAPI resets the paper counter to zero when staff replaces the media roll.
func (h *Handler) OperatorResetPaperAPI(c *gin.Context) {
	var body struct {
		Capacity int `json:"capacity"`
	}
	_ = c.ShouldBindJSON(&body)

	if err := settings.ResetPaperCounter(h.db, body.Capacity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cfgMap := settings.GetConfig(h.db, h.cfg.Port)
	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"paper_capacity":  cfgMap.PaperCapacity,
		"paper_used":      cfgMap.PaperUsed,
		"paper_remaining": cfgMap.PaperRemaining,
	})
}

// OperatorQuickVoucherAPI generates a single-use ticket code for fast hand-out to guests.
func (h *Handler) OperatorQuickVoucherAPI(c *gin.Context) {
	code := vouchers.GenerateRandomCode()
	v := vouchers.Voucher{
		Code:      code,
		IsUsed:    false,
		CreatedAt: time.Now(),
	}
	if err := h.db.Create(&v).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat kode voucher"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    code,
	})
}
