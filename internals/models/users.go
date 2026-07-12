package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName             string         `json:"full_name" gorm:"not null;size:100"`
	Email                string         `json:"email" gorm:"uniqueIndex;not null;size:255"`
	PhoneNumber          string         `json:"phone_number" gorm:"uniqueIndex;not null;size:20"`
	Password             string         `json:"-" gorm:"not null;size:255"`
	CareStreakDays       uint           `json:"care_streak_days" gorm:"default:0"`
	TotalTaskDone        uint           `json:"total_task_done" gorm:"default:0"`
	TotalJournalInjuries uint           `json:"total_journal_injuries" gorm:"default:0"`
	Notifications        []Notification `json:"notifications,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Plants               []Plant        `json:"plants,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Scans                []Scan         `json:"scans,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	AiChatSessions       []AiChatSession `json:"ai_chat_sessions,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
