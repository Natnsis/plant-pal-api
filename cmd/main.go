package main

import (
	"fmt"
	"log"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/models"
	"plantPal/internals/routes"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "plantPal/docs"
)

// @title           PlantPal API
// @version         1.0
// @description     API for PlantPal - a plant care and health monitoring application.
// @host            localhost:8080
// @BasePath        /
func main() {
	// config functions
	config.ConnectToDb()
	config.GetJwtSecret()

	// to auto migrate dbs
	models.MigrateDb()

	// initialze mux
	r := mux.NewRouter()
	// routes
	routes.AuthRoutes(r)

	// swagger docs endpoint
	r.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("./swagger.json"),
	))

	fmt.Println("server is running on port http://localhost:8080")
	fmt.Println("swagger docs available at http://localhost:8080/docs/")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("an error occured: ", err)
	}
}
