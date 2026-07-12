package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID                      uint       `json:"user_id" gorm:"index;not null"`
	NotificationEnabled         bool       `json:"notification_enabled" gorm:"default:true"`
	DailySummaryEnabled         bool       `json:"daily_summary_enabled" gorm:"default:false"`
	SoundAlertEnabled           bool       `json:"sound_alert_enabled" gorm:"default:true"`
	VibrationEnabled            bool       `json:"vibration_enabled" gorm:"default:true"`
	PreferredNotificationTime   time.Time  `json:"preferred_notification_time"`
	DefaultSnoozeDurationMinute uint       `json:"default_snooze_duration_minute" gorm:"default:15"`
	User                        *User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
