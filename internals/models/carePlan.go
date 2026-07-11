package models

import "gorm.io/gorm"

type CarePlan struct {
	gorm.Model
	PlantID               uint
	WateringFrequencyDays int
	WateringAmount        string
	WateringMethod        string
	WateringTips          string
	LightRequirment       string
	HumidityRequirment    string
}
