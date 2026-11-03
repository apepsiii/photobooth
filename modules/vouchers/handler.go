package vouchers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
	"mastergin/modules/settings"
)

// Handler handles voucher requests.
type Handler struct {
	db  *gorm.DB
	cfg *core.Config
}

// NewHandler creates a new Voucher handler.
func NewHandler(db *gorm.DB, cfg *core.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// AdminVouchersPage renders the vouchers management page in Meridian admin.
func (h *Handler) AdminVouchersPage(c *gin.Context) {
	list, err := ListAll(h.db)
	if err != nil {
		list = []Voucher{}
	}

	var totalCount int64
	var usedCount int64
	var availableCount int64

	h.db.Model(&Voucher{}).Count(&totalCount)
	h.db.Model(&Voucher{}).Where("is_used = ?", true).Count(&usedCount)
	availableCount = totalCount - usedCount

	voucherModeEnabled := settings.Get(h.db, "voucher_mode_enabled") == "true"

	c.HTML(http.StatusOK, "vouchers", core.PageData{
		Title:    "Tiket & Voucher Booth",
		UserName: middleware.CurrentUserName(c),
		Path:     "/admin/vouchers",
		Data: gin.H{
			"Vouchers":           list,
			"TotalCount":         totalCount,
			"UsedCount":          usedCount,
			"AvailableCount":     availableCount,
			"VoucherModeEnabled": voucherModeEnabled,
		},
	})
}

// Generate creates a batch of new vouchers.
func (h *Handler) Generate(c *gin.Context) {
	countStr := c.DefaultPostForm("count", "20")
	count, _ := strconv.Atoi(countStr)
	_, _ = GenerateBatch(h.db, count)
	c.Redirect(http.StatusFound, "/admin/vouchers")
}

// ToggleMode toggles voucher mode on/off.
func (h *Handler) ToggleMode(c *gin.Context) {
	current := settings.Get(h.db, "voucher_mode_enabled") == "true"
	newVal := "true"
	if current {
		newVal = "false"
	}
	_ = settings.Set(h.db, "voucher_mode_enabled", newVal)
	c.Redirect(http.StatusFound, "/admin/vouchers")
}

// Delete removes a voucher.
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	_ = Delete(h.db, uint(id))
	c.Redirect(http.StatusFound, "/admin/vouchers")
}

// ValidatePayload represents voucher validation request.
type ValidatePayload struct {
	Code string `json:"code"`
}

// APIValidate verifies a voucher code from the kiosk screen.
func (h *Handler) APIValidate(c *gin.Context) {
	var payload ValidatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"valid": false, "message": "Input tidak valid."})
		return
	}

	v, err := CheckValidity(h.db, payload.Code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"valid": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "code": v.Code})
}

