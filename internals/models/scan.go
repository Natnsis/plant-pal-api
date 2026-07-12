package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type Symptoms string

const (
	LeaveYellowSymptoms Symptoms = "yellow_leaves"
	BrownSpotSymptoms   Symptoms = "brown_spots"
	LeaveBrownSymptoms  Symptoms = "brown_steams"
	WeatheringSymptoms  Symptoms = "weathering_leaves"
)

type SymptomList []Symptoms

func (s SymptomList) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *SymptomList) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SymptomList: value is of type %T, not []byte", value)
	}
	return json.Unmarshal(bytes, s)
}

type Scan struct {
	gorm.Model
	UserID     uint
	PlantID    uint
	AnalysisID uint

	CapturedImageUrl string
	SelectedSymptoms SymptomList `gorm:"type:json"`
	AiOutputID       uint
}
