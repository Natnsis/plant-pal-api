package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityLog struct {
	gorm.Model
	PlantID      uint
	ActivityType string // enum watered, fertilized, repotted, photo node, milestone
	LoggedDate   time.Time
	Notes        string
	PhotoUrl     string
}
