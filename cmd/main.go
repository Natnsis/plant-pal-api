package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// config functions
	config.ConnectToDb()
	config.GetJwtSecret()
	config.GetCloudinaryConfig()
	config.GetGeminiAPIKey()

	// to auto migrate dbs
	models.MigrateDb()

	// initialze mux
	r := mux.NewRouter()
	r.Use(middlewares.CorsMiddleware)

	// health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")

	// swagger docs endpoint (public)
	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	})
	r.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	// auth routes (public)
	authRouter := r.PathPrefix("/").Subrouter()
	routes.AuthRoutes(authRouter)

	// protected routes (with JWT middleware)
	apiRouter := r.PathPrefix("/").Subrouter()
	apiRouter.Use(middlewares.AuthMiddleware)
	routes.PlantRoutes(apiRouter)
	routes.ScanRoutes(apiRouter)
	routes.CareRoutes(apiRouter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("server is running on port http://localhost:%s\n", port)
	fmt.Printf("swagger docs available at http://localhost:%s/docs/\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("an error occured: ", err)
	}
}
