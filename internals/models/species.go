package models

import "gorm.io/gorm"

type Species struct {
	gorm.Model
	CommonName      string
	ScientificName  string
	Family          string
	Origin          string
	DifficultyLevel int // enum
}
