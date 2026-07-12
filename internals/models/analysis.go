package models

import (
	"gorm.io/gorm"
)

type AnalysisType string

const (
	IdentificationAnalysisType AnalysisType = "identification"
	DiagnosisAnalysisType      AnalysisType = "diagnosis"
)

type AiAnalysisResult struct {
	gorm.Model
	ScanID             uint         `json:"scan_id" gorm:"uniqueIndex;not null"`
	AiModelVersion     string       `json:"ai_model_version" gorm:"size:50"`
	ConfidenceScore    float64      `json:"confidence_score" gorm:"check:confidence_score >= 0 AND confidence_score <= 1"`
	AnalysisType       AnalysisType `json:"analysis_type" gorm:"not null;size:30"`
	DetectedSymptoms   []string     `json:"detected_symptoms" gorm:"type:text[]"`
	PrimaryAssessment  string       `json:"primary_assessment" gorm:"type:text"`
	TreatmentPlanSteps string       `json:"treatment_plan_steps" gorm:"type:text"`
	MetadataPayload    string       `json:"metadata_payload" gorm:"type:text"`
	Scan               *Scan        `json:"scan,omitempty" gorm:"foreignKey:ScanID;constraint:OnDelete:CASCADE"`
}
