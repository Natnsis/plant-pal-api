package models

import "gorm.io/gorm"

type CarePlan struct {
	gorm.Model
	PlantID               uint   `json:"plant_id" gorm:"index;not null"`
	WateringFrequencyDays uint   `json:"watering_frequency_days" gorm:"not null;check:watering_frequency_days >= 1"`
	WateringAmount        string `json:"watering_amount" gorm:"size:50"`
	WateringMethod        string `json:"watering_method" gorm:"size:50"`
	WateringTips          string `json:"watering_tips" gorm:"type:text"`
	LightRequirement      string `json:"light_requirement" gorm:"size:100"`
	HumidityRequirement   string `json:"humidity_requirement" gorm:"size:100"`
	Plant                 *Plant `json:"plant,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
}
