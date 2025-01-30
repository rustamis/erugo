package db

import (
	"database/sql"
	"log"

	"github.com/DeanWard/erugo/models"
)

func GetUsers(database *sql.DB) []models.User {
	users := []models.User{}
	rows, err := database.Query("SELECT * FROM users")
	if err != nil {
		log.Println("Failed to get users", err)
		return []models.User{}
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		rows.Scan(&user.Username)
		users = append(users, user)
	}

	return users
}
