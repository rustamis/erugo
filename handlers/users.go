package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/models"
	"github.com/DeanWard/erugo/responses/json_response"
)

func GetUsersHandler(database *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleGetUsers(database, w)
	})
}

func CreateUserHandler(database *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCreateUser(database, w, r)
	})
}

func handleGetUsers(database *sql.DB, w http.ResponseWriter) {
	users, err := db.UserList(database)
	if err != nil {
		json_response.New(json_response.ErrorStatus, "Failed to get users", nil, http.StatusInternalServerError).Send(w)
		return
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	payload := map[string]interface{}{
		"metadata": map[string]interface{}{
			"total": len(users),
		},
		"users": userResponses,
	}

	json_response.New(json_response.SuccessStatus, "Users fetched successfully", payload, http.StatusOK).Send(w)
}

func handleCreateUser(database *sql.DB, w http.ResponseWriter, r *http.Request) {

	incomingUser := models.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&incomingUser)
	user := incomingUser.ToUser()

	if err != nil {
		json_response.New(json_response.ErrorStatus, "Failed to decode user", nil, http.StatusBadRequest).Send(w)
		return
	}

	errors := make(map[string]string)

	if incomingUser.Username == "" {
		errors["username"] = "Username is required"
	}

	if incomingUser.FullName == "" {
		errors["full_name"] = "Full name is required"
	}

	if incomingUser.Email == "" {
		errors["email"] = "Email is required"
	}

	if incomingUser.Password == "" {
		errors["password"] = "Password is required"
	}

	if len(errors) > 0 {
		payload := map[string]interface{}{
			"errors": errors,
		}
		json_response.New(json_response.ErrorStatus, "Validation errors", payload, http.StatusBadRequest).Send(w)
		return
	}

	createdUser, err := db.UserCreate(database, user)
	if err != nil {
		json_response.New(json_response.ErrorStatus, "Failed to create user", nil, http.StatusInternalServerError).Send(w)
		return
	}

	user = *createdUser
	payload := map[string]interface{}{
		"user": user.ToResponse(),
	}

	json_response.New(json_response.SuccessStatus, "User created successfully", payload, http.StatusCreated).Send(w)
}
