package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID                      uint
	NotificationEnabled         bool
	DailySummaryEnabled         bool
	SoundAlertEnabled           bool
	VibrationEnabled            bool
	PreferredNotificationTime   time.Time
	DefaultSnoozeDurationMinute int
}
