package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectToDb() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("unable to load env file")
		return
	}

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		log.Fatal("no database url found")
		return
	}

	conn, err := gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		log.Fatal("unable to connect to the db")
		return
	}

	Db = conn
}
