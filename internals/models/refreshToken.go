package models

import "time"

type RefreshToken struct {
	Token     string    `gorm:"uniqueIndex;not null" json:"-"`
	UserID    uint      `gorm:"not null;index" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"-"`
	Revoked   bool      `gorm:"default:false" json:"-"`
}
