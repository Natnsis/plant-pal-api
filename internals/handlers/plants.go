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

// ListPlants godoc
// @Summary      List all plants
// @Description  Get a list of all plants belonging to the authenticated user
// @Tags         plants
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Plant
// @Failure      500  {object}  response.ErrorResponse
// @Router       /plants [get]
func ListPlants(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var plants []models.Plant
	if result := config.Db.Where("user_id = ?", userID).Preload("Species").Preload("CarePlans").Find(&plants); result.Error != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch plants")
		return
	}

	response.JSON(w, http.StatusOK, plants)
}

// CreatePlant godoc
// @Summary      Create a new plant
// @Description  Add a new plant to the user's collection
// @Tags         plants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      CreatePlantRequest  true  "Plant details"
// @Success      201   {object}  models.Plant
// @Failure      400   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /plants [post]
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

// GetPlant godoc
// @Summary      Get a plant
// @Description  Get detailed information about a specific plant
// @Tags         plants
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Plant ID"
// @Success      200  {object}  models.Plant
// @Failure      404  {object}  response.ErrorResponse
// @Router       /plants/{id} [get]
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

// UpdatePlant godoc
// @Summary      Update a plant
// @Description  Update details of an existing plant
// @Tags         plants
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string              true  "Plant ID"
// @Param        body  body      UpdatePlantRequest  true  "Fields to update"
// @Success      200   {object}  models.Plant
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /plants/{id} [put]
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

// DeletePlant godoc
// @Summary      Delete a plant
// @Description  Remove a plant from the user's collection
// @Tags         plants
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Plant ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  response.ErrorResponse
// @Router       /plants/{id} [delete]
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
