package frames

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// PhotoSlot defines the coordinates and dimensions of an individual photo slot on a frame.
type PhotoSlot struct {
	Index        int `json:"index"`
	X            int `json:"x"`
	Y            int `json:"y"`
	Width        int `json:"w"`
	Height       int `json:"h"`
	BorderRadius int `json:"r,omitempty"`
}

// Frame defines a photo frame template.
type Frame struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:150;not null" json:"name"`
	Slug         string    `gorm:"size:150;uniqueIndex" json:"slug"`
	Category     string    `gorm:"size:100;default:'Umum'" json:"category"` // Pensi, Wisuda, Bazaar, dll
	Type         string    `gorm:"size:50;default:'strip_2x6'" json:"type"` // strip_2x6, grid_4cut, polaroid, custom
	Description  string    `gorm:"size:255" json:"description"`
	TotalShots   int       `gorm:"default:3" json:"total_shots"`
	Width        int       `gorm:"default:600" json:"width"`
	Height       int       `gorm:"default:1800" json:"height"`
	OverlayImage string    `gorm:"size:255" json:"overlay_image"` // Path to transparent PNG overlay
	SlotsJSON    string    `gorm:"type:text" json:"slots_json"`   // JSON array of PhotoSlot
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetSlots parses SlotsJSON into a slice of PhotoSlot.
func (f *Frame) GetSlots() []PhotoSlot {
	var slots []PhotoSlot
	if f.SlotsJSON != "" {
		_ = json.Unmarshal([]byte(f.SlotsJSON), &slots)
	}
	return slots
}

// AutoMigrate runs migrations for the Frame model.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Frame{})
}

