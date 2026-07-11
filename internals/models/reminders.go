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
	ReportTask    TaskType = "report"
)

type Reminder struct {
	gorm.Model
	PlantID       uint
	TaskType      TaskType
	ScheduledTime time.Time
	IsCompleted   bool
	CompletedAt   time.Time
	SnoozeCount   int
}
