package models

import (
	"time"

	"gorm.io/gorm"
)

type AiChatSession struct {
	gorm.Model
	UserID uint
	Status string // enum: active, archived
}
