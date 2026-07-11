package models

import (
	"time"

	"gorm.io/gorm"
)

type AiChat struct {
	gorm.Model
	SessionID   uint
	SenderType  string // enum user,ai
	MessageBody string
}
