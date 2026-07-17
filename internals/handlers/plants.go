package handlers

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"
	"plantPal/internals/response"

	"github.com/gorilla/mux"
)

type CreatePlantRequest struct {
	SpeciesID *uint  `json:"species_id"`
	Nickname  string `json:"nickname"`
	Location  string `json:"location"`
}

type UpdatePlantRequest struct {
	Nickname    *string `json:"nickname"`
	Location    *string `json:"location"`
	HealthScore *int    `json:"health_score"`
	Status      *string `json:"status"`
}

func ListPlants(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var plants []models.Plant
	if result := config.Db.Where("user_id = ?", userID).Preload("Species").Preload("CarePlans").Find(&plants); result.Error != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch plants")
		return
	}

	response.JSON(w, http.StatusOK, plants)
}

func CreatePlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var req CreatePlantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Nickname == "" {
		response.Error(w, http.StatusBadRequest, "nickname is required")
		return
	}

	plant := models.Plant{
		UserID:    userID,
		SpeciesID: req.SpeciesID,
		Nickname:  req.Nickname,
		Location:  req.Location,
		Status:    models.GoodPlantStatus,
	}

	if result := config.Db.Create(&plant); result.Error != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create plant")
		return
	}

	config.Db.Preload("Species").First(&plant, plant.ID)

	response.JSON(w, http.StatusCreated, plant)
}

func GetPlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).
		Preload("Species").Preload("CarePlans").Preload("Reminders").
		Preload("GrowthMetrics").Preload("ActivityLogs").Preload("Scans").
		First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	response.JSON(w, http.StatusOK, plant)
}

func UpdatePlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	var req UpdatePlantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.HealthScore != nil {
		updates["health_score"] = *req.HealthScore
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		config.Db.Model(&plant).Updates(updates)
	}

	config.Db.Preload("Species").First(&plant, plant.ID)

	response.JSON(w, http.StatusOK, plant)
}

func DeletePlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	config.Db.Delete(&plant)

	response.JSON(w, http.StatusOK, map[string]string{"message": "plant deleted successfully"})
}
