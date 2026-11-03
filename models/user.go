package models

import (
	"time"
)

// User represents an authenticated application user.
// This model is auto-migrated by GORM on application startup.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"default:user;not null" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName overrides the default table name (would be "users" anyway, but
// made explicit for AI readability).
func (User) TableName() string {
	return "users"
}
