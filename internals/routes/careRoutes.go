package routes

import (
	"plantPal/internals/handlers"

	"github.com/gorilla/mux"
)

func CareRoutes(r *mux.Router) {
	r.HandleFunc("/plants/{id}/care-plan", handlers.GetCarePlan).Methods("GET")
	r.HandleFunc("/plants/{id}/care-plan", handlers.UpdateCarePlan).Methods("PUT")

	r.HandleFunc("/plants/{id}/reminders", handlers.GetPlantReminders).Methods("GET")
	r.HandleFunc("/plants/{id}/activities", handlers.GetActivities).Methods("GET")
	r.HandleFunc("/plants/{id}/activities", handlers.CreateActivity).Methods("POST")
	r.HandleFunc("/plants/{id}/growth", handlers.GetGrowthMetrics).Methods("GET")
	r.HandleFunc("/plants/{id}/growth", handlers.CreateGrowthMetric).Methods("POST")

	r.HandleFunc("/reminders/today", handlers.GetTodayReminders).Methods("GET")
	r.HandleFunc("/reminders/{id}", handlers.UpdateReminder).Methods("PUT")

	r.HandleFunc("/notifications", handlers.GetNotificationSettings).Methods("GET")
	r.HandleFunc("/notifications", handlers.UpdateNotificationSettings).Methods("PUT")
}
