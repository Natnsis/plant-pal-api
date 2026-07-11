package models

import "gorm.io/gorm"

type Plant struct {
	gorm.Model
	UserID      uint
	SpeciesID   uint
	Nickname    string
	Location    string
	HealthScore int    //0-100
	Status      string // enum
}
