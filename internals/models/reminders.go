package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskType string

const (
	WaterTask     TaskType = "water"
	FertilizeTask TaskType = "fertilize"
	MistTask      TaskType = "mist"
	RotateTask    TaskType = "rotate"
	RepotTask     TaskType = "repot"
)

type Reminder struct {
	gorm.Model
	PlantID       uint      `json:"plant_id" gorm:"index;not null"`
	TaskType      TaskType  `json:"task_type" gorm:"not null;size:20"`
	ScheduledTime time.Time `json:"scheduled_time" gorm:"not null"`
	IsCompleted   bool      `json:"is_completed" gorm:"default:false"`
	CompletedAt   time.Time `json:"completed_at"`
	SnoozeCount   uint      `json:"snooze_count" gorm:"default:0"`
	Plant         *Plant    `json:"plant,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
}
