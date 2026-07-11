package models

import "gorm.io/gorm"

type DifficultyType string

const (
	EasyDifficult   DifficultyType = "easy"
	MediumDifficult DifficultyType = "medium"
	HardDifficult   DifficultyType = "hard"
)

type Species struct {
	gorm.Model
	CommonName      string
	ScientificName  string
	Family          string
	Origin          string
	DifficultyLevel DifficultyType
}
