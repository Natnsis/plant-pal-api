package models

import (
	"time"

	"gorm.io/gorm"
)

type SenderType string

const (
	SenderTypeUser SenderType = "user"
	SenderTypeAi   SenderType = "ai"
)

type AiChat struct {
	gorm.Model
	SessionID   uint
	SenderType  SenderType
	MessageBody string
}
