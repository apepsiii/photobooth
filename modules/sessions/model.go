package sessions

import (
	"time"

	"gorm.io/gorm"

	"mastergin/modules/frames"
)

// Session represents a single photobooth event session.
type Session struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	SessionCode     string         `gorm:"size:20;uniqueIndex;not null" json:"session_code"`
	FrameID         uint           `json:"frame_id"`
	Frame           frames.Frame   `gorm:"foreignKey:FrameID" json:"frame"`
	CompositePath   string         `gorm:"size:255" json:"composite_path"`
	FilterApplied   string         `gorm:"size:50;default:'normal'" json:"filter_applied"`
	PrintedCount    int            `gorm:"default:0" json:"printed_count"`
	DownloadedCount int            `gorm:"default:0" json:"downloaded_count"`
	Status          string         `gorm:"size:30;default:'completed'" json:"status"`
	Photos          []SessionPhoto `gorm:"foreignKey:SessionID" json:"photos"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// SessionPhoto represents an individual raw camera capture from a session.
type SessionPhoto struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"index;not null" json:"session_id"`
	ShotOrder int       `json:"shot_order"`
	PhotoPath string    `gorm:"size:255" json:"photo_path"`
	CreatedAt time.Time `json:"created_at"`
}

// AutoMigrate runs migrations for sessions models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Session{}, &SessionPhoto{})
}

