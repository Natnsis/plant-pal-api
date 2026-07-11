package models

import (
	"time"

	"gorm.io/gorm"
)

type Reminder struct {
	gorm.Model
	PlantID       uint
	TaskType      string // enum=> water, fertilize, mist, rotate, report
	ScheduledTime time.Time
	IsCompleted   bool
	CompletedAt   time.Time // nullable
	SnoozeCount   int
}
