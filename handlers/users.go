package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/models"
	"github.com/gorilla/mux"
)

// Response is a standardized API response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// sendResponse is a helper function to send JSON responses
func sendResponse(w http.ResponseWriter, status string, message string, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If we fail to encode the response, we've already set the header status,
		// so we'll just log the error and return
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetUsersHandler returns all users
func GetUsersHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.UserList(database)
		if err != nil {
			sendResponse(w, StatusError, "Failed to fetch users", nil, http.StatusInternalServerError)
			return
		}

		// Convert to response objects
		userResponses := make([]models.UserResponse, len(users))
		for i, user := range users {
			userResponses[i] = user.ToResponse()
		}

		data := map[string]interface{}{
			"users": userResponses,
			"metadata": map[string]interface{}{
				"total": len(users),
			},
		}

		sendResponse(w, StatusSuccess, "Users retrieved successfully", data, http.StatusOK)
	}
}

// CreateUserHandler creates a new user
func CreateUserHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createReq models.UserCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			sendResponse(w, StatusError, "Invalid request format", nil, http.StatusBadRequest)
			return
		}

		// Validate the request
		if err := createReq.Validate(); err != nil {
			sendResponse(w, StatusError, err.Error(), nil, http.StatusBadRequest)
			return
		}

		// Convert request to user model
		user := createReq.ToUser()

		// Create the user in the database
		createdUser, err := db.UserCreate(database, user)
		if err != nil {
			sendResponse(w, StatusError, "Failed to create user", nil, http.StatusInternalServerError)
			return
		}

		sendResponse(w, StatusSuccess, "User created successfully",
			map[string]interface{}{"user": createdUser.ToResponse()},
			http.StatusCreated)
	}
}

// UpdateUserHandler updates an existing user
func UpdateUserHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate the user ID
		userID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			sendResponse(w, StatusError, "Invalid user ID", nil, http.StatusBadRequest)
			return
		}

		// Parse the update request
		var updateReq models.UserUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			sendResponse(w, StatusError, "Invalid request format", nil, http.StatusBadRequest)
			return
		}

		// Get the existing user
		existingUser, err := db.UserByID(database, userID)
		if err != nil {
			if errors.Is(err, db.ErrUserNotFound) {
				sendResponse(w, StatusError, "User not found", nil, http.StatusNotFound)
				return
			}
			sendResponse(w, StatusError, "Failed to fetch user", nil, http.StatusInternalServerError)
			return
		}

		// Apply updates to the existing user
		if updateReq.Username != nil {
			existingUser.Username = *updateReq.Username
		}
		if updateReq.Password != nil {
			if err := db.UserSetPassword(database, userID, *updateReq.Password); err != nil {
				sendResponse(w, StatusError, "Failed to update password", nil, http.StatusInternalServerError)
				return
			}
		}
		if updateReq.Admin != nil {
			existingUser.Admin = *updateReq.Admin
		}
		if updateReq.FullName != nil {
			existingUser.FullName = sql.NullString{
				String: *updateReq.FullName,
				Valid:  *updateReq.FullName != "",
			}
		}
		if updateReq.Email != nil {
			existingUser.Email = sql.NullString{
				String: *updateReq.Email,
				Valid:  *updateReq.Email != "",
			}
		}
		if updateReq.Active != nil {
			existingUser.Active = *updateReq.Active
		}

		// Update the user in the database
		updatedUser, err := db.UserUpdate(database, *existingUser)
		if err != nil {
			sendResponse(w, StatusError, "Failed to update user", nil, http.StatusInternalServerError)
			return
		}

		sendResponse(w, StatusSuccess, "User updated successfully",
			map[string]interface{}{"user": updatedUser.ToResponse()},
			http.StatusOK)
	}
}

// DeleteUserHandler removes a user
func DeleteUserHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate the user ID
		userID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			sendResponse(w, StatusError, "Invalid user ID", nil, http.StatusBadRequest)
			return
		}

		// Delete the user
		err = db.UserDelete(database, userID)
		if err != nil {
			if errors.Is(err, db.ErrUserNotFound) {
				sendResponse(w, StatusError, "User not found", nil, http.StatusNotFound)
				return
			}
			sendResponse(w, StatusError, "Failed to delete user", nil, http.StatusInternalServerError)
			return
		}

		sendResponse(w, StatusSuccess, "User deleted successfully", nil, http.StatusOK)
	}
}
