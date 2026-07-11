package models

import (
	"time"

	"gorm.io/gorm"
)

type GrowthRateStatus string

const (
	SlowGrowthRate     GrowthRateStauts = "slow"
	ModerateGrowthRate GrowthRateStauts = "moderate"
	FastGrowthRate     GrowthRateStauts = "fast"
)

type GrowthMetric struct {
	gorm.Model
	PlantID          uint
	RecordedDate     time.Time
	HeightCm         Float
	GrowthRateStauts GrowthRateStauts
}
