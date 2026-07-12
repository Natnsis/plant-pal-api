package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityType string

const (
	ActivityWatered    ActivityType = "watered"
	ActivityFertilized ActivityType = "fertilized"
	ActivityRepotted   ActivityType = "repotted"
	ActivityPhotoNode  ActivityType = "photo_node"
	ActivityMilestone  ActivityType = "milestone"
)

type ActivityLog struct {
	gorm.Model
	PlantID      uint         `json:"plant_id" gorm:"index;not null"`
	ActivityType ActivityType `json:"activity_type" gorm:"not null;size:20"`
	LoggedDate   time.Time    `json:"logged_date" gorm:"not null"`
	Notes        string       `json:"notes" gorm:"type:text"`
	PhotoUrl     string       `json:"photo_url" gorm:"size:500"`
	Plant        *Plant       `json:"plant,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
}
