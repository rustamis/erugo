package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/DeanWard/erugo/config"
	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/models"
	"github.com/DeanWard/erugo/responses"
	"github.com/DeanWard/erugo/validation"
)

func GetSettingsByGroupHandler(database *sql.DB) http.HandlerFunc {
	validator := validation.NewSettingValidator(database)

	return func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")

		// Validate the group parameter
		if validationErrors := validator.ValidateGroup(group); validationErrors.HasErrors() {
			responses.SendResponse(w, responses.StatusError, "Validation failed",
				map[string]interface{}{"errors": validationErrors},
				http.StatusBadRequest)
			return
		}

		settings, err := db.SettingsByGroup(database, group)
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to fetch settings", nil, http.StatusInternalServerError)
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

		responses.SendResponse(w, responses.StatusSuccess, "Settings retrieved successfully", data, http.StatusOK)
	}
}

func GetSettingByIdHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if id == "" {
			responses.SendResponse(w, responses.StatusError, "Setting ID is required", nil, http.StatusBadRequest)
			return
		}

		setting, err := db.SettingById(database, id)
		if err != nil {
			if errors.Is(err, db.ErrSettingNotFound) {
				responses.SendResponse(w, responses.StatusError, "Setting not found", nil, http.StatusNotFound)
				return
			}
			responses.SendResponse(w, responses.StatusError, "Failed to fetch setting", nil, http.StatusInternalServerError)
			return
		}

		responses.SendResponse(w, responses.StatusSuccess, "Setting retrieved successfully", setting.ToResponse(), http.StatusOK)
	}
}

func SetSettingsByIdHandler(database *sql.DB) http.HandlerFunc {
	validator := validation.NewSettingValidator(database)

	return func(w http.ResponseWriter, r *http.Request) {

		var settingsReq []models.SettingRequest
		if err := json.NewDecoder(r.Body).Decode(&settingsReq); err != nil {
			responses.SendResponse(w, responses.StatusError, "Invalid request format", nil, http.StatusBadRequest)
			return
		}

		// Validate all settings first
		allErrors := make(models.ValidationErrors)
		for i, settingReq := range settingsReq {
			if errors := validator.ValidateSet(settingReq.Id, settingReq.Value, settingReq.SettingGroup); len(errors) > 0 {
				// Add index prefix to error keys to identify which setting had the error
				for field, msgs := range errors {
					key := fmt.Sprintf("%d.%s", i, field)
					allErrors[key] = msgs
				}
			}
		}

		if len(allErrors) > 0 {
			responses.SendResponse(w, responses.StatusError, "Validation failed",
				map[string]interface{}{"errors": allErrors},
				http.StatusBadRequest)
			return
		}

		// Process each setting
		updatedSettings := make([]models.SettingResponse, 0, len(settingsReq))
		for _, settingReq := range settingsReq {
			// Get the group from the request or try to get existing setting's group
			group := settingReq.SettingGroup
			if group == "" {
				// If no group provided, try to get existing setting's group
				existing, err := db.SettingById(database, settingReq.Id)
				if err != nil {
					if errors.Is(err, db.ErrSettingNotFound) {
						responses.SendResponse(w, responses.StatusError, fmt.Sprintf("Setting %s not found", settingReq.Id), nil, http.StatusNotFound)
						return
					}
					responses.SendResponse(w, responses.StatusError, "Failed to check existing setting", nil, http.StatusInternalServerError)
					return
				}
				group = existing.SettingGroup
			}

			// Update or create the setting
			err := db.SettingSetById(database, settingReq.Id, settingReq.Value, group)
			if err != nil {
				if errors.Is(err, db.ErrSettingNotFound) {
					responses.SendResponse(w, responses.StatusError, fmt.Sprintf("Setting %s not found", settingReq.Id), nil, http.StatusNotFound)
					return
				}
				responses.SendResponse(w, responses.StatusError, "Failed to set setting", nil, http.StatusInternalServerError)
				return
			}

			// Fetch the updated setting
			setting, err := db.SettingById(database, settingReq.Id)
			if err != nil {
				responses.SendResponse(w, responses.StatusError, "Failed to fetch updated setting", nil, http.StatusInternalServerError)
				return
			}
			updatedSettings = append(updatedSettings, setting.ToResponse())
		}

		responses.SendResponse(w, responses.StatusSuccess, "Settings updated successfully", updatedSettings, http.StatusOK)
	}
}

func SetLogoHandler() http.HandlerFunc {
	privateDir := config.AppConfig.PrivateDataPath
	return func(w http.ResponseWriter, r *http.Request) {
		//we're going to receive a file from the request. We'll store the file in the private directory as logo.<ext>
		file, _, err := r.FormFile("logo")
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to get logo file", nil, http.StatusBadRequest)
			return
		}
		defer file.Close()

		//create the private directory if it doesn't exist
		os.MkdirAll(privateDir, 0755)

		//create the logo file
		logoFile, err := os.Create(filepath.Join(privateDir, "logo.png"))
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to create logo file", nil, http.StatusInternalServerError)
			return
		}
		defer logoFile.Close()

		//copy the file to the logo file
		io.Copy(logoFile, file)

		responses.SendResponse(w, responses.StatusSuccess, "Logo updated successfully", nil, http.StatusOK)
	}
}

func GetLogoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//read the logo file
		logoFile, err := os.ReadFile(filepath.Join("./private", "logo.png"))
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to read logo file", nil, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(logoFile)
	}
}
