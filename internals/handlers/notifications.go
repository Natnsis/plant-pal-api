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

// GetNotificationSettings godoc
// @Summary      Get notification settings
// @Description  Retrieve the notification preferences for the authenticated user
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.Notification
// @Failure      500  {object}  response.ErrorResponse
// @Router       /notifications [get]
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

// UpdateNotificationSettings godoc
// @Summary      Update notification settings
// @Description  Update notification preferences for the authenticated user
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      UpdateNotificationRequest  true  "Settings to update"
// @Success      200   {object}  models.Notification
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /notifications [put]
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
