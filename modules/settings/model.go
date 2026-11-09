package settings

import (
	"time"

	"gorm.io/gorm"
)

// Setting represents a key-value application configuration in SQLite.
type Setting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Label     string    `gorm:"size:255" json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ConfigMap holds typed settings for the photobooth application.
type ConfigMap struct {
	AppName            string `json:"app_name"`
	EventName          string `json:"event_name"`
	SchoolName         string `json:"school_name"`
	CountdownSeconds   int    `json:"countdown_seconds"`
	PrintEnabled       bool   `json:"print_enabled"`
	AutoPrint          bool   `json:"auto_print"`
	CopiesCount        int    `json:"copies_count"`
	PrinterName        string `json:"printer_name"`
	KioskPin           string `json:"kiosk_pin"`
	LocalServerIP      string `json:"local_server_ip"`
	Port               string `json:"port"`
	Domain             string `json:"domain"`
	BaseURL            string `json:"base_url"`
	VoucherModeEnabled bool   `json:"voucher_mode_enabled"`
	PaymentModeEnabled bool   `json:"payment_mode_enabled"`
	PriceSingle        int    `json:"price_single"`
	PriceDouble        int    `json:"price_double"`
	QRISStaticCode     string `json:"qris_static_code"`
	PaperCapacity      int    `json:"paper_capacity"`
	PaperUsed          int    `json:"paper_used"`
	PaperRemaining     int    `json:"paper_remaining"`
}

// AutoMigrate migrates the Setting model.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Setting{})
}

