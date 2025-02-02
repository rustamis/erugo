package db

import (
	"database/sql"
	"log"

	"github.com/DeanWard/erugo/models"
	"golang.org/x/crypto/bcrypt"
)

func UserList(database *sql.DB) ([]models.User, error) {
	users := []models.User{}
	rows, err := database.Query("SELECT id, username, admin, full_name, email, must_change_pw, created_at, updated_at, active FROM users")
	if err != nil {
		log.Println("Failed to get users", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Admin, &user.FullName, &user.Email, &user.MustChangePassword, &user.CreatedAt, &user.UpdatedAt, &user.Active); err != nil {
			log.Println("Failed to scan user", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func UserByName(database *sql.DB, username string) *models.User {
	log.Printf("Getting user by username: %s", username)
	row := database.QueryRow("SELECT id, username, password_hash, admin, full_name, email, must_change_pw, created_at, updated_at, active FROM users WHERE username = ?", username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Admin, &user.FullName, &user.Email, &user.MustChangePassword, &user.CreatedAt, &user.UpdatedAt, &user.Active)
	if err != nil {
		return nil
	}
	return &user
}

func UserByID(database *sql.DB, id int) *models.User {
	log.Printf("Getting user by id: %d", id)
	row := database.QueryRow("SELECT id, username, admin, full_name, email, must_change_pw, created_at, updated_at, active FROM users WHERE id = ?", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Admin, &user.FullName, &user.Email, &user.MustChangePassword, &user.CreatedAt, &user.UpdatedAt, &user.Active)
	if err != nil {
		return nil
	}
	return &user
}

func UserCreate(database *sql.DB, user models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_, err = database.Exec("INSERT INTO users (username, password_hash, admin, full_name, email, must_change_pw, created_at, updated_at, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", user.Username, hashedPassword, user.Admin, user.FullName, user.Email, user.MustChangePassword, user.CreatedAt, user.UpdatedAt, user.Active)
	if err != nil {
		return nil, err
	}
	//grab and return the created user
	user.Password = ""
	return UserByName(database, user.Username), nil
}

func UserUpdate(database *sql.DB, user models.User) (*models.User, error) {
	_, err := database.Exec("UPDATE users SET username = ?, admin = ?, full_name = ?, email = ?, active = ? WHERE id = ?", user.Username, user.Admin, user.FullName, user.Email, user.Active, user.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserSetPassword(database *sql.DB, user models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_, err = database.Exec("UPDATE users SET password_hash = ? WHERE id = ?", hashedPassword, user.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserDelete(database *sql.DB, user models.User) error {
	_, err := database.Exec("DELETE FROM users WHERE id = ?", user.ID)
	if err != nil {
		return err
	}
	return nil
}
