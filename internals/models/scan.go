package models

import (
	"gorm.io/gorm"
)

type Symptoms string

const (
	YellowLeavesSymptom Symptoms = "yellow_leaves"
	BrownSpotsSymptom   Symptoms = "brown_spots"
	BrownStemsSymptom   Symptoms = "brown_stems"
	WitheringLeavesSymptom Symptoms = "withering_leaves"
)

type Scan struct {
	gorm.Model
	UserID           uint       `json:"user_id" gorm:"index;not null"`
	PlantID          uint       `json:"plant_id" gorm:"index;not null"`
	CapturedImageUrl string     `json:"captured_image_url" gorm:"size:500"`
	SelectedSymptoms []Symptoms `json:"selected_symptoms" gorm:"type:text[]"`
	User             *User      `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Plant            *Plant     `json:"plant,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
}
