package handlers

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"
	"plantPal/internals/response"
)

type UpdateNotificationRequest struct {
	NotificationEnabled         *bool   `json:"notification_enabled"`
	DailySummaryEnabled         *bool   `json:"daily_summary_enabled"`
	SoundAlertEnabled           *bool   `json:"sound_alert_enabled"`
	VibrationEnabled            *bool   `json:"vibration_enabled"`
	PreferredNotificationTime   *string `json:"preferred_notification_time"`
	DefaultSnoozeDurationMinute *uint   `json:"default_snooze_duration_minute"`
}

func GetNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var notification models.Notification
	if result := config.Db.Where("user_id = ?", userID).First(&notification); result.Error != nil {
		notification = models.Notification{
			UserID:              userID,
			NotificationEnabled: true,
			DailySummaryEnabled: false,
			SoundAlertEnabled:   true,
			VibrationEnabled:    true,
		}
		config.Db.Create(&notification)
	}

	response.JSON(w, http.StatusOK, notification)
}

func UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var notification models.Notification
	if result := config.Db.Where("user_id = ?", userID).First(&notification); result.Error != nil {
		response.Error(w, http.StatusNotFound, "settings not found")
		return
	}

	var req UpdateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updates := map[string]interface{}{}
	if req.NotificationEnabled != nil {
		updates["notification_enabled"] = *req.NotificationEnabled
	}
	if req.DailySummaryEnabled != nil {
		updates["daily_summary_enabled"] = *req.DailySummaryEnabled
	}
	if req.SoundAlertEnabled != nil {
		updates["sound_alert_enabled"] = *req.SoundAlertEnabled
	}
	if req.VibrationEnabled != nil {
		updates["vibration_enabled"] = *req.VibrationEnabled
	}
	if req.DefaultSnoozeDurationMinute != nil {
		updates["default_snooze_duration_minute"] = *req.DefaultSnoozeDurationMinute
	}
	if req.PreferredNotificationTime != nil {
		updates["preferred_notification_time"] = *req.PreferredNotificationTime
	}

	if len(updates) > 0 {
		config.Db.Model(&notification).Updates(updates)
	}

	response.JSON(w, http.StatusOK, notification)
}
