package db

import (
	"database/sql"
	"log"

	"github.com/DeanWard/erugo/models"
)

func GetUsers(database *sql.DB) ([]models.User, error) {
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
