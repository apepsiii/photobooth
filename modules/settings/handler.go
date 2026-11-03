package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/middleware"
)

// Handler handles settings HTTP requests.
type Handler struct {
	db  *gorm.DB
	cfg *core.Config
}

// NewHandler creates a new Handler instance.
func NewHandler(db *gorm.DB, cfg *core.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// SettingsPage renders the settings page in the Meridian admin layout.
func (h *Handler) SettingsPage(c *gin.Context) {
	config := GetConfig(h.db, h.cfg.Port)

	c.HTML(http.StatusOK, "settings_app", core.PageData{
		Title:    "Pengaturan Kiosk & Printer",
		UserName: middleware.CurrentUserName(c),
		Path:     "/admin/settings",
		Data: gin.H{
			"Config":          config,
			"DetectedLocalIP": DetectLocalIP(),
		},
	})
}

// UpdateSettings processes the settings update form.
func (h *Handler) UpdateSettings(c *gin.Context) {
	_ = Set(h.db, "app_name", c.PostForm("app_name"))
	_ = Set(h.db, "event_name", c.PostForm("event_name"))
	_ = Set(h.db, "school_name", c.PostForm("school_name"))
	_ = Set(h.db, "countdown_seconds", c.PostForm("countdown_seconds"))
	_ = Set(h.db, "copies_count", c.PostForm("copies_count"))
	_ = Set(h.db, "printer_name", c.PostForm("printer_name"))
	_ = Set(h.db, "kiosk_pin", c.PostForm("kiosk_pin"))
	_ = Set(h.db, "local_server_ip", c.PostForm("local_server_ip"))

	if c.PostForm("print_enabled") == "on" || c.PostForm("print_enabled") == "true" {
		_ = Set(h.db, "print_enabled", "true")
	} else {
		_ = Set(h.db, "print_enabled", "false")
	}

	if c.PostForm("auto_print") == "on" || c.PostForm("auto_print") == "true" {
		_ = Set(h.db, "auto_print", "true")
	} else {
		_ = Set(h.db, "auto_print", "false")
	}

	config := GetConfig(h.db, h.cfg.Port)

	c.HTML(http.StatusOK, "settings_app", core.PageData{
		Title:     "Pengaturan Kiosk & Printer",
		UserName:  middleware.CurrentUserName(c),
		Path:      "/admin/settings",
		Flash:     "Pengaturan berhasil disimpan!",
		FlashType: "success",
		Data: gin.H{
			"Config":          config,
			"DetectedLocalIP": DetectLocalIP(),
		},
	})
}

// GetConfigJSON returns the active config as JSON (for kiosk frontend).
func (h *Handler) GetConfigJSON(c *gin.Context) {
	config := GetConfig(h.db, h.cfg.Port)
	c.JSON(http.StatusOK, config)
}

