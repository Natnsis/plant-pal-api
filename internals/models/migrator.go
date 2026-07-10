package models

import (
	"fmt"

	"plantPal/internals/config"
)

func MigrateDb() {
	err := config.Db.AutoMigrate(
		&User{},
	)
	if err != nil {
		// if ther is error on migration
		fmt.Printf("database migration failed: %v", err)
		return
	}

	// if everyting is fine
	fmt.Printf("models are migrated")
}
