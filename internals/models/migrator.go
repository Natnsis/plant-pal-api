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
		fmt.Println("user is not migrated")
	} else {
		fmt.Println("unable to migrate db")
	}
}
