package routes

import (
	"plantPal/internals/auth"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/register", auth.Register).Methods("POST")
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.HandleFunc("/refresh", auth.Refresh).Methods("POST")
	r.HandleFunc("/logout", auth.Logout).Methods("POST")
	r.HandleFunc("/auth/google", auth.LoginWithGoogle).Methods("POST")
}
