package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type AnalysisType string

const (
	IdentificationAnalysisType AnalysisType = "identification"
	DiagnosisAnalysisType      AnalysisType = "diagnosis"
)

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringList: value is of type %T, not []byte", value)
	}
	return json.Unmarshal(bytes, s)
}

type AiAnalysisResult struct {
	gorm.Model
	ScanID             uint         `json:"scan_id" gorm:"uniqueIndex;not null"`
	AiModelVersion     string       `json:"ai_model_version" gorm:"size:50"`
	ConfidenceScore    float64      `json:"confidence_score" gorm:"check:confidence_score >= 0 AND confidence_score <= 1"`
	AnalysisType       AnalysisType `json:"analysis_type" gorm:"not null;size:30"`
	DetectedSymptoms   StringList   `json:"detected_symptoms" gorm:"type:json"`
	PrimaryAssessment  string       `json:"primary_assessment" gorm:"type:text"`
	TreatmentPlanSteps string       `json:"treatment_plan_steps" gorm:"type:text"`
	MetadataPayload    string       `json:"metadata_payload" gorm:"type:text"`
	Scan               *Scan        `json:"scan,omitempty" gorm:"foreignKey:ScanID;constraint:OnDelete:CASCADE"`
}
