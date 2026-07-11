package models

import (
	"time"

	"gorm.io/gorm"
)

type AiAnalysisResult struct {
	gorm.Model
	ScanID             uint
	AiModelVersion     string
	ConfidenceScore    float64
	AnalysisType       string // enum identificatin, diagnosis
	DetectedSymptoms   []string
	PrimaryAssessment  string
	TreatmentPlanSteps string
	MetadataPayload    string
}
