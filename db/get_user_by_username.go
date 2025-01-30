package db

import (
	"database/sql"
	"log"

	"github.com/DeanWard/erugo/models"
)

func GetUserByUsername(database *sql.DB, username string) *models.User {
	log.Printf("Getting user by username: %s", username)
	row := database.QueryRow("SELECT id, username, password_hash, admin FROM users WHERE username = ?", username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Admin)
	if err != nil {
		return nil
	}
	return &user
}
