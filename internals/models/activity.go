package models

import "gorm.io/gorm"

type ActivityLog struct {
	gorm.Model
	Title        string
	Description  string
	Status       string
	ActivityType string
}
