package models

import (
	"time"

	"gorm.io/gorm"
)

type AiSessionStatus string

const (
	ActiveSession   AiSessionStatus = "active"
	ArchivedSession AiSessionStatus = "archived"
)

type AiSession struct {
	gorm.Model
	UserID uint
	Status AiSessionStatus
}
