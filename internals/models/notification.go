package models

import "gorm.io/gorm"

type Reminders struct {
	gorm.Model
	PlantID  uint
	Title    string
	SubTitle string
	IsRead   bool
}
