package models

import "gorm.io/gorm"

type PlantStatus string

const (
	GoodPlantStatus        PlantStatus = "good"
	NeedsAttentionStatus   PlantStatus = "needs_attention"
)

type Plant struct {
	gorm.Model
	UserID         uint           `json:"user_id" gorm:"index;not null"`
	SpeciesID      *uint          `json:"species_id" gorm:"index"`
	Nickname       string         `json:"nickname" gorm:"not null;size:100"`
	Location       string         `json:"location" gorm:"size:100"`
	HealthScore    int            `json:"health_score" gorm:"default:0;check:health_score >= 0 AND health_score <= 100"`
	Status         PlantStatus    `json:"status" gorm:"not null;size:20;default:'good'"`
	User           *User          `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Species        *Species       `json:"species,omitempty" gorm:"foreignKey:SpeciesID;constraint:OnDelete:SET NULL"`
	CarePlans      []CarePlan     `json:"care_plans,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
	Reminders      []Reminder     `json:"reminders,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
	GrowthMetrics  []GrowthMetric `json:"growth_metrics,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
	ActivityLogs   []ActivityLog  `json:"activity_logs,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
	Scans          []Scan         `json:"scans,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
}
