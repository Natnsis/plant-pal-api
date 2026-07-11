package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityType string

const (
	StatusWatered    ActivityType = "watered"
	StatusFertilized ActivityType = "fertilized"
	StatusRepotted   ActivityType = "repotted"
	StatusPhotoNode  ActivityType = "photo_node"
	StatusMilestore  ActivityType = "milestone"
)

type ActivityLog struct {
	gorm.Model
	PlantID      uint
	ActivityType ActivityType
	LoggedDate   time.Time
	Notes        string
	PhotoUrl     string
}
