package routes

import (
	"plantPal/internals/handlers"

	"github.com/gorilla/mux"
)

func PlantRoutes(r *mux.Router) {
	r.HandleFunc("/plants", handlers.ListPlants).Methods("GET")
	r.HandleFunc("/plants", handlers.CreatePlant).Methods("POST")
	r.HandleFunc("/plants/{id}", handlers.GetPlant).Methods("GET")
	r.HandleFunc("/plants/{id}", handlers.UpdatePlant).Methods("PUT")
	r.HandleFunc("/plants/{id}", handlers.DeletePlant).Methods("DELETE")
}
