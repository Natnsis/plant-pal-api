package models

import (
	"time"

	"gorm.io/gorm"
)

type GrowthRateStatus string

const (
	SlowGrowthRate     GrowthRateStatus = "slow"
	ModerateGrowthRate GrowthRateStatus = "moderate"
	FastGrowthRate     GrowthRateStatus = "fast"
)

type GrowthMetric struct {
	gorm.Model
	PlantID          uint             `json:"plant_id" gorm:"index;not null"`
	RecordedDate     time.Time        `json:"recorded_date" gorm:"not null"`
	HeightCm         float64          `json:"height_cm" gorm:"check:height_cm >= 0"`
	GrowthRateStatus GrowthRateStatus `json:"growth_rate_status" gorm:"not null;size:20"`
	Plant            *Plant           `json:"plant,omitempty" gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
}
