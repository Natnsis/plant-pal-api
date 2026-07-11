package models

import (
	"time"

	"gorm.io/gorm"
)

type GrowthMetric struct {
	gorm.Model
	PlantID          uint
	RecordedDate     time.Time
	HeightCm         Float
	GrowthRateStauts string // slow, moderate, fast
}
