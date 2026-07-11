package models

import (
	"time"

	"gorm.io/gorm"
)

type AnalysisType string

const (
	IdentificationAnalysisType AnalysisType = "indentification"
	DiagnosisAnalysisType      AnalysisType = "diagnosis"
)

type AiAnalysisResult struct {
	gorm.Model
	ScanID             uint
	AiModelVersion     string
	ConfidenceScore    float64
	AnalysisType       AnalysisType
	DetectedSymptoms   []string
	PrimaryAssessment  string
	TreatmentPlanSteps string
	MetadataPayload    string
}
