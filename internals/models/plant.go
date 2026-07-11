package models

import "gorm.io/gorm"

type PlantStatus string

const (
	GoodPlantStatus PlantStatus = "good"
	BadStatus       PlantStatus = "needs_attention"
)

type Plant struct {
	gorm.Model
	UserID      uint
	SpeciesID   uint
	Nickname    string
	Location    string
	HealthScore int //0-100
	Status      PlantStatus
}
