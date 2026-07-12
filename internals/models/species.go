package models

import "gorm.io/gorm"

type DifficultyType string

const (
	EasyDifficulty   DifficultyType = "easy"
	MediumDifficulty DifficultyType = "medium"
	HardDifficulty   DifficultyType = "hard"
)

type Species struct {
	gorm.Model
	CommonName      string         `json:"common_name" gorm:"not null;size:100"`
	ScientificName  string         `json:"scientific_name" gorm:"not null;size:150"`
	Family          string         `json:"family" gorm:"size:100"`
	Origin          string         `json:"origin" gorm:"size:100"`
	DifficultyLevel DifficultyType `json:"difficulty_level" gorm:"not null;size:20"`
	Plants          []Plant        `json:"plants,omitempty" gorm:"foreignKey:SpeciesID;constraint:OnDelete:SET NULL"`
}
