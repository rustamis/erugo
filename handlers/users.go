package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/middleware"
	"github.com/DeanWard/erugo/models"
	"github.com/DeanWard/erugo/responses"
	"github.com/DeanWard/erugo/validation"
	"github.com/gorilla/mux"
)

// Response is a standardized API response structure

// GetUsersHandler returns all users
func GetUsersHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.UserList(database)
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to fetch users", nil, http.StatusInternalServerError)
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

		responses.SendResponse(w, responses.StatusSuccess, "Users retrieved successfully", data, http.StatusOK)
	}
}

// CreateUserHandler creates a new user

func CreateUserHandler(database *sql.DB) http.HandlerFunc {
	validator := validation.NewUserValidator(database)

	return func(w http.ResponseWriter, r *http.Request) {
		var createReq models.UserCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			responses.SendResponse(w, responses.StatusError, "Invalid request format", nil, http.StatusBadRequest)
			return
		}

		// Validate including uniqueness checks
		if validationErrors := validator.ValidateCreate(&createReq); validationErrors.HasErrors() {
			responses.SendResponse(w, responses.StatusError, "Validation failed",
				map[string]interface{}{"errors": validationErrors},
				http.StatusBadRequest)
			return
		}

		// Convert request to user model
		user := createReq.ToUser()

		// Create the user in the database
		createdUser, err := db.UserCreate(database, user)
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to create user", nil, http.StatusInternalServerError)
			return
		}

		responses.SendResponse(w, responses.StatusSuccess, "User created successfully",
			map[string]interface{}{"user": createdUser.ToResponse()},
			http.StatusCreated)
	}
}

// UpdateUserHandler updates an existing user
func UpdateUserHandler(database *sql.DB) http.HandlerFunc {
	validator := validation.NewUserValidator(database)

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate the user ID from URL
		userID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Invalid user ID", nil, http.StatusBadRequest)
			return
		}

		// Parse the update request body
		var updateReq models.UserUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			responses.SendResponse(w, responses.StatusError, "Invalid request format", nil, http.StatusBadRequest)
			return
		}

		// Get the existing user
		existingUser, err := db.UserByID(database, userID)
		if err != nil {
			if errors.Is(err, db.ErrUserNotFound) {
				responses.SendResponse(w, responses.StatusError, "User not found", nil, http.StatusNotFound)
				return
			}
			responses.SendResponse(w, responses.StatusError, "Failed to fetch user", nil, http.StatusInternalServerError)
			return
		}

		// Validate updates including uniqueness checks
		if validationErrors := validator.ValidateUpdate(existingUser, &updateReq, false, false); validationErrors.HasErrors() {
			responses.SendResponse(w, responses.StatusError, "Validation failed",
				map[string]interface{}{"errors": validationErrors},
				http.StatusBadRequest)
			return
		}

		// Update password if provided
		if updateReq.Password != nil {
			if err := db.UserSetPassword(database, userID, *updateReq.Password); err != nil {
				responses.SendResponse(w, responses.StatusError, "Failed to update password", nil, http.StatusInternalServerError)
				return
			}
		}

		// Apply updates to the existing user
		if updateReq.Username != nil {
			existingUser.Username = *updateReq.Username
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
		if updateReq.MustChangePassword != nil {
			existingUser.MustChangePassword = *updateReq.MustChangePassword
		}

		if updateReq.Active != nil {
			existingUser.Active = *updateReq.Active
		}

		// Update the user in the database
		updatedUser, err := db.UserUpdate(database, *existingUser)
		if err != nil {
			if errors.Is(err, db.ErrUserNotFound) {
				responses.SendResponse(w, responses.StatusError, "User not found", nil, http.StatusNotFound)
				return
			}
			responses.SendResponse(w, responses.StatusError, "Failed to update user", nil, http.StatusInternalServerError)
			return
		}

		responses.SendResponse(w, responses.StatusSuccess, "User updated successfully",
			map[string]interface{}{"user": updatedUser.ToResponse()},
			http.StatusOK)
	}
}

// DeleteUserHandler removes a user
func DeleteUserHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//this user id is in the context
		CurrentUserID := r.Context().Value(middleware.ContextKey("userID")).(int)

		// Parse and validate the user ID
		userID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Invalid user ID", nil, http.StatusBadRequest)
			return
		}

		// Check if the user is the current user
		if userID == CurrentUserID {
			responses.SendResponse(w, responses.StatusError, "You cannot delete yourself", nil, http.StatusForbidden)
			return
		}

		// Delete the user
		err = db.UserDelete(database, userID)
		if err != nil {
			if errors.Is(err, db.ErrUserNotFound) {
				responses.SendResponse(w, responses.StatusError, "User not found", nil, http.StatusNotFound)
				return
			}
			responses.SendResponse(w, responses.StatusError, "Failed to delete user", nil, http.StatusInternalServerError)
			return
		}

		responses.SendResponse(w, responses.StatusSuccess, "User deleted successfully", nil, http.StatusOK)
	}
}

// GetMyProfileHandler returns the current user
func GetMyProfileHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		CurrentUserID := r.Context().Value(middleware.ContextKey("userID")).(int)
		user, err := db.UserByID(database, CurrentUserID)
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to fetch user", nil, http.StatusInternalServerError)
			return
		}

		responses.SendResponse(w, responses.StatusSuccess, "User retrieved successfully", user.ToResponse(), http.StatusOK)
	}
}

// UpdateMyProfileHandler updates the current user
func UpdateMyProfileHandler(database *sql.DB) http.HandlerFunc {
	validator := validation.NewUserValidator(database)
	return func(w http.ResponseWriter, r *http.Request) {
		CurrentUserID := r.Context().Value(middleware.ContextKey("userID")).(int)
		user, err := db.UserByID(database, CurrentUserID)
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to fetch user", nil, http.StatusInternalServerError)
			return
		}

		// Parse the update request body
		var updateReq models.UserUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			responses.SendResponse(w, responses.StatusError, "Invalid request format", nil, http.StatusBadRequest)
			return
		}

		// Validate updates including uniqueness checks
		if validationErrors := validator.ValidateUpdate(user, &updateReq, true, true); validationErrors.HasErrors() {
			responses.SendResponse(w, responses.StatusError, "Validation failed",
				map[string]interface{}{"errors": validationErrors},
				http.StatusBadRequest)
			return
		}

		if updateReq.Username != nil {
			user.Username = *updateReq.Username
		}
		if updateReq.FullName != nil {
			user.FullName = sql.NullString{
				String: *updateReq.FullName,
				Valid:  *updateReq.FullName != "",
			}
		}
		if updateReq.Email != nil {
			user.Email = sql.NullString{
				String: *updateReq.Email,
				Valid:  *updateReq.Email != "",
			}
		}
		if updateReq.Password != nil {

			if err := db.UserSetPassword(database, CurrentUserID, *updateReq.Password); err != nil {
				responses.SendResponse(w, responses.StatusError, "Failed to update password", nil, http.StatusInternalServerError)
				return
			}
		}

		// Update the user in the database
		updatedUser, err := db.UserUpdate(database, *user)
		if err != nil {
			responses.SendResponse(w, responses.StatusError, "Failed to update user", nil, http.StatusInternalServerError)
			return
		}

		responses.SendResponse(w, responses.StatusSuccess, "User updated successfully", updatedUser.ToResponse(), http.StatusOK)
	}
}
