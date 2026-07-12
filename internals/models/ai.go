package models

import "gorm.io/gorm"

type ChatStatus string

const (
	ChatStatusActive   ChatStatus = "active"
	ChatStatusArchived ChatStatus = "archived"
)

type AiChatSession struct {
	gorm.Model
	UserID  uint       `json:"user_id" gorm:"index;not null"`
	Status  ChatStatus `json:"status" gorm:"not null;size:20;default:'active'"`
	User    *User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	AiChats []AiChat   `json:"ai_chats,omitempty" gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
}
