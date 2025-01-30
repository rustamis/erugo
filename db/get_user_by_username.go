package db

import (
	"database/sql"
	"log"

	"github.com/DeanWard/erugo/models"
)

func GetUserByUsername(database *sql.DB, username string) *models.User {
	log.Printf("Getting user by username: %s", username)
	row := database.QueryRow("SELECT username, password_hash FROM users WHERE username = ?", username)
	var user models.User
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		return nil
	}
	return &user
}
