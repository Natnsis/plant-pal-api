package config

import (
	"fmt"
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
		fmt.Println("no .env file found, using environment variables")
	} else {
		fmt.Println("env is loaded")
	}

	db_url := os.Getenv("LOCAL_DB_URL")
	if db_url == "" {
		db_url = os.Getenv("DATABASE_URL")
	}
	if db_url == "" {
		log.Fatal("no database url found")
		return
	}
	fmt.Println("db url is fetched")

	conn, err := gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		log.Fatal("unable to connect to the db")
		return
	}
	fmt.Println("connected to db")

	Db = conn
}
