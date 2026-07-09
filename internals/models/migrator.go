package models

import (
	"plantPal/internals/config"
)

func MigrateDb() {
	config.Db.AutoMigrate(
		&User{},
	)
}
