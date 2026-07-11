package models

import (
	"time"

	"gorm.io/gorm"
)

type Scan struct {
	gorm.Model
	UserID           uint
	PlantID          uint
	ScanType         string // enum: identification, diagnosis
	CapturedImageUrl string
	SelectedSymptoms []string // yellow leaves, brown spots..
	AiOutputID       uint
}
