package main

import (
	"fmt"
	"log"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/models"
)

func main() {
	// connection
	config.ConnectToDb()
	// to auto migrate dbs
	models.MigrateDb()

	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		log.Fatal("an error occured")
	}

	fmt.Println("server is running on port :8080")
}
