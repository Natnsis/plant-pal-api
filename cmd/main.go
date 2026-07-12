package main

import (
	"fmt"
	"log"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/models"
	"plantPal/internals/routes"

	"github.com/gorilla/mux"
)

func main() {
	// connection
	config.ConnectToDb()
	// to auto migrate dbs
	models.MigrateDb()

	// initialze mux
	r := mux.NewRouter()
	// routes
	routes.AuthRoutes(r)

	fmt.Println("server is running on port http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("an error occured: ", err)
	}
}
