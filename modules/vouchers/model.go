package vouchers

import (
	"time"

	"gorm.io/gorm"
)

// Voucher represents a prepaid token or ticket for school photobooth sessions.
type Voucher struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Code        string     `gorm:"size:30;uniqueIndex;not null" json:"code"`
	IsUsed      bool       `gorm:"default:false" json:"is_used"`
	UsedAt      *time.Time `json:"used_at,omitempty"`
	SessionCode string     `gorm:"size:30" json:"session_code,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// AutoMigrate runs migrations for Voucher model.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Voucher{})
}

