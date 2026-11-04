package kiosk

import (
	"time"

	"gorm.io/gorm"
)

// PaymentOrder represents a self-service kiosk payment transaction.
type PaymentOrder struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	OrderID     string     `gorm:"size:50;uniqueIndex;not null" json:"order_id"`
	PackageType string     `gorm:"size:20;not null" json:"package_type"` // "single", "double"
	Amount      int        `gorm:"not null" json:"amount"`
	Copies      int        `gorm:"default:1" json:"copies"`
	Status      string     `gorm:"size:20;default:'pending'" json:"status"` // "pending", "paid", "expired"
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AutoMigrate runs migrations for kiosk payment models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&PaymentOrder{})
}

