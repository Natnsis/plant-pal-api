package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"
	"plantPal/internals/response"
	"plantPal/internals/services"

	"github.com/gorilla/mux"
)

// ScanPlant godoc
// @Summary      Identify a plant
// @Description  Upload a plant image for AI-powered identification. Returns a preview for the user to confirm before saving.
// @Tags         scan
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        image  formData  file  true  "Plant image"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  response.ErrorResponse
// @Failure      429    {object}  response.ErrorResponse
// @Failure      500    {object}  response.ErrorResponse
// @Router       /scan [post]
func ScanPlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var todayCount int64
	today := time.Now().Truncate(24 * time.Hour)
	config.Db.Model(&models.Scan{}).Where("user_id = ? AND created_at >= ?", userID, today).Count(&todayCount)
	if todayCount >= 5 {
		response.Error(w, http.StatusTooManyRequests, "daily scan limit reached (5 per day)")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	imageURL, err := services.UploadImage(file, fileHeader)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("failed to upload image: %s", err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	identification, err := services.IdentifyPlant(ctx, imageURL)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("failed to analyze plant: %s", err.Error()))
		return
	}

	identificationJSON, _ := json.Marshal(identification)

	scan := models.Scan{
		UserID:                    userID,
		CapturedImageUrl:          imageURL,
		ConfidenceScore:           identification.ConfidenceScore,
		JsonIdentificationPayload: string(identificationJSON),
		Retake:                    identification.ConfidenceScore < 0.7,
	}
	config.Db.Create(&scan)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"scan_id":            scan.ID,
		"retake":             scan.Retake,
		"confidence_score":   identification.ConfidenceScore,
		"identification":     identification,
		"captured_image_url": imageURL,
	})
}

// ConfirmScanRequest godoc
// @Summary      Confirm a scan and create plant
// @Description  Confirm that the identified plant is correct, creating species, plant, care plan, and reminders
// @Tags         scan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string            true  "Scan ID"
// @Param        body  body      ConfirmScanRequest  true  "Confirmation payload"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /scan/{id}/confirm [post]
func ConfirmScan(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	scanID := mux.Vars(r)["id"]

	var scan models.Scan
	if result := config.Db.Where("id = ? AND user_id = ?", scanID, userID).First(&scan); result.Error != nil {
		response.Error(w, http.StatusNotFound, "scan not found")
		return
	}

	if scan.PlantID != 0 {
		response.Error(w, http.StatusBadRequest, "scan already confirmed")
		return
	}

	var req ConfirmScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var identification services.PlantIdentification
	if err := json.Unmarshal([]byte(scan.JsonIdentificationPayload), &identification); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to parse identification data")
		return
	}

	nickname := req.Nickname
	if nickname == "" {
		nickname = identification.CommonName
	}

	species := models.Species{
		CommonName:      identification.CommonName,
		ScientificName:  identification.ScientificName,
		Family:          identification.Family,
		Origin:          identification.Origin,
		DifficultyLevel: models.MediumDifficulty,
	}
	config.Db.Create(&species)

	plant := models.Plant{
		UserID:      userID,
		SpeciesID:   &species.ID,
		Nickname:    nickname,
		Location:    req.Location,
		HealthScore: int(identification.ConfidenceScore * 100),
		Status:      models.GoodPlantStatus,
	}
	config.Db.Create(&plant)

	analysis := models.AiAnalysisResult{
		ScanID:             scan.ID,
		AiModelVersion:     "gemini-2.0-flash",
		ConfidenceScore:    identification.ConfidenceScore,
		AnalysisType:       models.IdentificationAnalysisType,
		DetectedSymptoms:   models.StringList(identification.DetectedSymptoms),
		PrimaryAssessment:  identification.PrimaryAssessment,
		TreatmentPlanSteps: joinStrings(identification.TreatmentSteps),
	}
	config.Db.Create(&analysis)

	carePlan := models.CarePlan{
		PlantID:               plant.ID,
		WateringFrequencyDays: identification.CareRecommendations.WateringFrequencyDays,
		WateringAmount:        identification.CareRecommendations.WateringAmount,
		WateringMethod:        identification.CareRecommendations.WateringMethod,
		WateringTips:          identification.CareRecommendations.WateringTips,
		LightRequirement:      identification.CareRecommendations.LightRequirement,
		HumidityRequirement:   identification.CareRecommendations.HumidityRequirement,
	}
	config.Db.Create(&carePlan)

	now := time.Now()
	createReminderIfNotExists(plant.ID, models.WaterTask, now.AddDate(0, 0, int(identification.CareRecommendations.WateringFrequencyDays)))
	createReminderIfNotExists(plant.ID, models.FertilizeTask, now.AddDate(0, 1, 0))
	createReminderIfNotExists(plant.ID, models.RotateTask, now.AddDate(0, 0, 7))

	scan.PlantID = plant.ID
	scan.AnalysisID = analysis.ID
	config.Db.Save(&scan)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"scan_id":            scan.ID,
		"plant":              plant,
		"species":            species,
		"analysis":           analysis,
		"care_plan":          carePlan,
		"captured_image_url": scan.CapturedImageUrl,
	})
}

// GetScan godoc
// @Summary      Get scan details
// @Description  Retrieve details of a specific scan including analysis results
// @Tags         scan
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Scan ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  response.ErrorResponse
// @Router       /scan/{id} [get]
func GetScan(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	scanID := mux.Vars(r)["id"]

	var scan models.Scan
	if result := config.Db.Where("id = ? AND user_id = ?", scanID, userID).
		Preload("Plant").Preload("Plant.Species").First(&scan); result.Error != nil {
		response.Error(w, http.StatusNotFound, "scan not found")
		return
	}

	var analysis models.AiAnalysisResult
	config.Db.Where("scan_id = ?", scan.ID).First(&analysis)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"scan":     scan,
		"analysis": analysis,
	})
}

type ConfirmScanRequest struct {
	Nickname string `json:"nickname"`
	Location string `json:"location"`
}

func createReminderIfNotExists(plantID uint, taskType models.TaskType, scheduledTime time.Time) {
	var count int64
	config.Db.Model(&models.Reminder{}).
		Where("plant_id = ? AND task_type = ? AND is_completed = false", plantID, taskType).
		Count(&count)
	if count == 0 {
		reminder := models.Reminder{
			PlantID:       plantID,
			TaskType:      taskType,
			ScheduledTime: scheduledTime,
		}
		config.Db.Create(&reminder)
	}
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += "\n"
		}
		result += s
	}
	return result
}
