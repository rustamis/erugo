package db

import (
	"database/sql"
	"log"

	"github.com/DeanWard/erugo/models"
)

func UserList(database *sql.DB) ([]models.User, error) {
	users := []models.User{}
	rows, err := database.Query("SELECT * FROM users")
	if err != nil {
		log.Println("Failed to get users", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Admin); err != nil {
			log.Println("Failed to scan user", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func UserByName(database *sql.DB, username string) *models.User {
	log.Printf("Getting user by username: %s", username)
	row := database.QueryRow("SELECT id, username, password_hash, admin FROM users WHERE username = ?", username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Admin)
	if err != nil {
		return nil
	}
	return &user
}
