package models

import "gorm.io/gorm"

type SenderType string

const (
	SenderTypeUser SenderType = "user"
	SenderTypeAI   SenderType = "ai"
)

type AiChat struct {
	gorm.Model
	SessionID   uint        `json:"session_id" gorm:"index;not null"`
	SenderType  SenderType  `json:"sender_type" gorm:"not null;size:10"`
	MessageBody string      `json:"message_body" gorm:"type:text;not null"`
	AiChatSession *AiChatSession `json:"ai_chat_session,omitempty" gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
}
