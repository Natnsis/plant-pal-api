package models

import (
	"time"

	"gorm.io/gorm"
)

type Symptoms string

const (
	LeaveSymptoms     Symptoms = "yellow_leaves"
	BrownSpotSymptoms Symptoms = "brown_spots"
	LeaveSymptoms     Symptoms = "brown_steams"
	LeaveSymptoms     Symptoms = "weathering_leaves"
)

type Scan struct {
	gorm.Model
	UserID           uint
	PlantID          uint
	AnalysisID       uint
	CapturedImageUrl string
	SelectedSymptoms []Symptoms
	AiOutputID       uint
}
