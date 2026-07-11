package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatStatus string

const (
	StatusActive   ChatStatus = "active"
	StatusArchived ChatStatus = "archived"
)

type AiChatSession struct {
	gorm.Model
	UserID uint
	Status ChatStatus
}
