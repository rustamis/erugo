package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DeanWard/erugo/models"
	"golang.org/x/crypto/bcrypt"
)

// Common database errors
var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrDBOperation  = fmt.Errorf("database operation failed")
)

// UserList retrieves all users from the database
func UserList(db *sql.DB) ([]models.User, error) {
	query := `
		SELECT id, username, admin, full_name, email, 
		       must_change_pw, created_at, updated_at, active 
		FROM users
		ORDER BY id`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Admin,
			&user.FullName, &user.Email, &user.MustChangePassword,
			&user.CreatedAt, &user.UpdatedAt, &user.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	return users, nil
}

// UserByName retrieves a user by their username
func UserByName(db *sql.DB, username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, admin, full_name, 
		       email, must_change_pw, created_at, updated_at, active 
		FROM users 
		WHERE username = ?`

	var user models.User
	err := db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Admin,
		&user.FullName, &user.Email, &user.MustChangePassword,
		&user.CreatedAt, &user.UpdatedAt, &user.Active,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	return &user, nil
}

// UserByID retrieves a user by their ID
func UserByID(db *sql.DB, id int) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, admin, full_name, 
		       email, must_change_pw, created_at, updated_at, active 
		FROM users 
		WHERE id = ?`

	var user models.User
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Admin,
		&user.FullName, &user.Email, &user.MustChangePassword,
		&user.CreatedAt, &user.UpdatedAt, &user.Active,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	return &user, nil
}

// UserCreate creates a new user in the database
func UserCreate(db *sql.DB, user models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.PasswordHash),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		INSERT INTO users (
			username, password_hash, admin, full_name, email,
			must_change_pw, created_at, updated_at, active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query,
		user.Username, hashedPassword, user.Admin,
		user.FullName, user.Email, user.MustChangePassword,
		user.CreatedAt, user.UpdatedAt, user.Active,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	return UserByID(db, int(id))
}

// UserUpdate updates an existing user's information
func UserUpdate(db *sql.DB, user models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users 
		SET username = ?, admin = ?, full_name = ?, 
		    email = ?, active = ?, updated_at = ? 
		WHERE id = ?`

	result, err := db.Exec(query,
		user.Username, user.Admin, user.FullName,
		user.Email, user.Active, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}
	if rows == 0 {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

// UserSetPassword updates a user's password
func UserSetPassword(db *sql.DB, userID int, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
		UPDATE users 
		SET password_hash = ?, updated_at = ? 
		WHERE id = ?`

	result, err := db.Exec(query, hashedPassword, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBOperation, err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UserDelete removes a user from the database
func UserDelete(db *sql.DB, userID int) error {
	result, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBOperation, err)
	}
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
