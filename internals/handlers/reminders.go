package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"

	"github.com/gorilla/mux"
)

// GetPlantReminders godoc
// @Summary      Get reminders for a plant
// @Description  Get all reminders for a specific plant
// @Tags         reminders
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Success      200 {array} models.Reminder
// @Failure      401 {string} string "unauthorized"
// @Router       /plants/{id}/reminders [get]
func GetPlantReminders(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	var reminders []models.Reminder
	config.Db.Where("plant_id = ?", plant.ID).Order("scheduled_time asc").Find(&reminders)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}

// GetTodayReminders godoc
// @Summary      Get today's reminders
// @Description  Get all due reminders for the authenticated user today
// @Tags         reminders
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Reminder
// @Failure      401 {string} string "unauthorized"
// @Router       /reminders/today [get]
func GetTodayReminders(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	var reminders []models.Reminder
	config.Db.Joins("JOIN plants ON plants.id = reminders.plant_id").
		Where("plants.user_id = ? AND reminders.scheduled_time >= ? AND reminders.scheduled_time < ? AND reminders.is_completed = false",
			userID, startOfDay, endOfDay).
		Order("reminders.scheduled_time asc").
		Find(&reminders)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}

type UpdateReminderRequest struct {
	IsCompleted *bool `json:"is_completed"`
	Snooze      bool  `json:"snooze"`
}

// UpdateReminder godoc
// @Summary      Update a reminder
// @Description  Complete or snooze a reminder
// @Tags         reminders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Reminder ID"
// @Param        body body UpdateReminderRequest true "Update payload"
// @Success      200 {object} models.Reminder
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "reminder not found"
// @Router       /reminders/{id} [put]
func UpdateReminder(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	reminderID := mux.Vars(r)["id"]

	var reminder models.Reminder
	if result := config.Db.Joins("JOIN plants ON plants.id = reminders.plant_id").
		Where("reminders.id = ? AND plants.user_id = ?", reminderID, userID).
		First(&reminder); result.Error != nil {
		http.Error(w, "reminder not found", http.StatusNotFound)
		return
	}

	var req UpdateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.IsCompleted != nil && *req.IsCompleted {
		reminder.IsCompleted = true
		reminder.CompletedAt = time.Now()
		config.Db.Save(&reminder)

		// Update user task count
		var plant models.Plant
		config.Db.Where("id = ?", reminder.PlantID).First(&plant)
		config.Db.Model(&models.User{}).Where("id = ?", plant.UserID).
			UpdateColumn("total_task_done", 1)
	} else if req.Snooze {
		reminder.SnoozeCount++
		reminder.ScheduledTime = reminder.ScheduledTime.Add(15 * time.Minute)
		config.Db.Save(&reminder)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminder)
}
