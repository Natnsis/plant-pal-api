package models

import (
	"log"

	"plantPal/internals/config"
)

func MigrateDb() {
	err := config.Db.AutoMigrate(
		&User{},
		&RefreshToken{},
		&Species{},
		&Plant{},
		&Scan{},
		&AiAnalysisResult{},
		&CarePlan{},
		&GrowthMetric{},
		&ActivityLog{},
		&Reminder{},
		&Notification{},
		&AiChatSession{},
		&AiChat{},
	)
	if err != nil {
		log.Fatal("failed to auto migrate database: ", err)
	}
}
