package handlers

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"

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

// ListPlants godoc
// @Summary      List user's plants
// @Description  Get all plants belonging to the authenticated user
// @Tags         plants
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array}  models.Plant
// @Failure      401 {string} string "unauthorized"
// @Router       /plants [get]
func ListPlants(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var plants []models.Plant
	if result := config.Db.Where("user_id = ?", userID).Preload("Species").Preload("CarePlans").Find(&plants); result.Error != nil {
		http.Error(w, "failed to fetch plants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plants)
}

// CreatePlant godoc
// @Summary      Create a plant
// @Description  Manually create a new plant entry
// @Tags         plants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body CreatePlantRequest true "Plant payload"
// @Success      201 {object} models.Plant
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Router       /plants [post]
func CreatePlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var req CreatePlantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Nickname == "" {
		http.Error(w, "nickname is required", http.StatusBadRequest)
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
		http.Error(w, "failed to create plant", http.StatusInternalServerError)
		return
	}

	config.Db.Preload("Species").First(&plant, plant.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(plant)
}

// GetPlant godoc
// @Summary      Get plant details
// @Description  Get a single plant with all related data
// @Tags         plants
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Success      200 {object} models.Plant
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "plant not found"
// @Router       /plants/{id} [get]
func GetPlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).
		Preload("Species").Preload("CarePlans").Preload("Reminders").
		Preload("GrowthMetrics").Preload("ActivityLogs").Preload("Scans").
		First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plant)
}

// UpdatePlant godoc
// @Summary      Update a plant
// @Description  Update plant details
// @Tags         plants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Param        body body UpdatePlantRequest true "Update payload"
// @Success      200 {object} models.Plant
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "plant not found"
// @Router       /plants/{id} [put]
func UpdatePlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	var req UpdatePlantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plant)
}

// DeletePlant godoc
// @Summary      Delete a plant
// @Description  Delete a plant and all its related data
// @Tags         plants
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Success      200 {object} map[string]string
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "plant not found"
// @Router       /plants/{id} [delete]
func DeletePlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	config.Db.Delete(&plant)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "plant deleted successfully"})
}
