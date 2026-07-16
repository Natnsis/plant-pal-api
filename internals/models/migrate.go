package models

import (
	"log"

	"plantPal/internals/config"
)

func MigrateDb() {
	config.Db.Exec("DROP TABLE IF EXISTS scans CASCADE")
	config.Db.Exec("DROP TABLE IF EXISTS ai_analysis_results CASCADE")

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
