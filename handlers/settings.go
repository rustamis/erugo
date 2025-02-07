package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/models"
	"github.com/DeanWard/erugo/validation"
)

func GetSettingsByGroupHandler(database *sql.DB) http.HandlerFunc {
	validator := validation.NewSettingValidator(database)

	return func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")

		// Validate the group parameter
		if validationErrors := validator.ValidateGroup(group); validationErrors.HasErrors() {
			sendResponse(w, StatusError, "Validation failed",
				map[string]interface{}{"errors": validationErrors},
				http.StatusBadRequest)
			return
		}

		settings, err := db.SettingsByGroup(database, group)
		if err != nil {
			sendResponse(w, StatusError, "Failed to fetch settings", nil, http.StatusInternalServerError)
			return
		}

		// Convert to response objects
		settingResponses := make([]models.SettingResponse, len(settings))
		for i, setting := range settings {
			settingResponses[i] = setting.ToResponse()
		}

		data := map[string]interface{}{
			"settings": settingResponses,
			"metadata": map[string]interface{}{
				"total": len(settings),
			},
		}

		sendResponse(w, StatusSuccess, "Settings retrieved successfully", data, http.StatusOK)
	}
}

func GetSettingByIdHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if id == "" {
			sendResponse(w, StatusError, "Setting ID is required", nil, http.StatusBadRequest)
			return
		}

		setting, err := db.SettingById(database, id)
		if err != nil {
			if errors.Is(err, db.ErrSettingNotFound) {
				sendResponse(w, StatusError, "Setting not found", nil, http.StatusNotFound)
				return
			}
			sendResponse(w, StatusError, "Failed to fetch setting", nil, http.StatusInternalServerError)
			return
		}

		sendResponse(w, StatusSuccess, "Setting retrieved successfully", setting.ToResponse(), http.StatusOK)
	}
}

func SetSettingByIdHandler(database *sql.DB) http.HandlerFunc {
	validator := validation.NewSettingValidator(database)

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var settingReq models.SettingRequest
		if err := json.NewDecoder(r.Body).Decode(&settingReq); err != nil {
			sendResponse(w, StatusError, "Invalid request format", nil, http.StatusBadRequest)
			return
		}

		// Validate the input
		if validationErrors := validator.ValidateSet(settingReq.Id, settingReq.Value, settingReq.SettingGroup); validationErrors.HasErrors() {
			sendResponse(w, StatusError, "Validation failed",
				map[string]interface{}{"errors": validationErrors},
				http.StatusBadRequest)
			return
		}

		// Get the group from the request or try to get existing setting's group
		group := settingReq.SettingGroup
		if group == "" {
			// If no group provided, try to get existing setting's group
			existing, err := db.SettingById(database, settingReq.Id)
			if err != nil && !errors.Is(err, db.ErrSettingNotFound) {
				sendResponse(w, StatusError, "Failed to check existing setting", nil, http.StatusInternalServerError)
				return
			}
			if existing != nil {
				group = existing.SettingGroup
			} else {
				sendResponse(w, StatusError, "Setting group is required for new settings", nil, http.StatusBadRequest)
				return
			}
		}

		err := db.SettingSetById(database, settingReq.Id, settingReq.Value, group)
		if err != nil {
			if errors.Is(err, db.ErrSettingNotFound) {
				sendResponse(w, StatusError, "Setting not found", nil, http.StatusNotFound)
				return
			}
			sendResponse(w, StatusError, "Failed to set setting", nil, http.StatusInternalServerError)
			return
		}

		// Fetch the updated setting to return in response
		setting, err := db.SettingById(database, settingReq.Id)
		if err != nil {
			sendResponse(w, StatusError, "Failed to fetch updated setting", nil, http.StatusInternalServerError)
			return
		}

		sendResponse(w, StatusSuccess, "Setting updated successfully", setting.ToResponse(), http.StatusOK)
	}
}
