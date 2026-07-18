package routes

import (
	"plantPal/internals/handlers"

	"github.com/gorilla/mux"
)

func ScanRoutes(r *mux.Router) {
	r.HandleFunc("/scan", handlers.ScanPlant).Methods("POST")
	r.HandleFunc("/scan/{id}", handlers.GetScan).Methods("GET")
	r.HandleFunc("/scan/{id}/confirm", handlers.ConfirmScan).Methods("POST")

	r.HandleFunc("/diagnosis", handlers.StartDiagnosis).Methods("POST")
	r.HandleFunc("/diagnosis/{session_id}", handlers.GetDiagnosisChat).Methods("GET")
	r.HandleFunc("/diagnosis/{session_id}/chat", handlers.SendChatMessage).Methods("POST")
}
