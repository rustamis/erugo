package handlers

import (
	"database/sql"
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
