package models

import (
	"database/sql"
	"errors"
	"regexp"
	"time"
)

var (
	// Common validation errors
	ErrMissingRequired = errors.New("required field is missing")
	ErrInvalidFormat   = errors.New("invalid format")

	// Email validation regex
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// User represents the database model
type User struct {
	ID                 int            `json:"id"`
	Username           string         `json:"username"`
	PasswordHash       string         `json:"-"` // Never expose in JSON
	Admin              bool           `json:"admin"`
	FullName           sql.NullString `json:"full_name"`
	Email              sql.NullString `json:"email"`
	MustChangePassword bool           `json:"must_change_password"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Active             bool           `json:"active"`
}

// UserResponse is the standardized API response structure
type UserResponse struct {
	ID                 int       `json:"id"`
	Username           string    `json:"username"`
	Admin              bool      `json:"admin"`
	FullName           string    `json:"full_name,omitempty"`
	Email              string    `json:"email,omitempty"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Active             bool      `json:"active"`
}

// UserCreateRequest represents the expected input for user creation
type UserCreateRequest struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	Admin              bool   `json:"admin"`
	FullName           string `json:"full_name"`
	Email              string `json:"email"`
	MustChangePassword bool   `json:"must_change_password"`
}

// UserUpdateRequest represents the expected input for user updates
type UserUpdateRequest struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Admin    *bool   `json:"admin,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Active   *bool   `json:"active,omitempty"`
}

// Validate performs validation on the create request
func (r *UserCreateRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username: " + ErrMissingRequired.Error())
	}
	if r.Password == "" {
		return errors.New("password: " + ErrMissingRequired.Error())
	}
	if r.Email != "" && !emailRegex.MatchString(r.Email) {
		return errors.New("email: " + ErrInvalidFormat.Error())
	}
	return nil
}

// ToUser converts a creation request to a User model
func (r *UserCreateRequest) ToUser() User {
	return User{
		Username:     r.Username,
		PasswordHash: "", // Will be set by the handler
		Admin:        r.Admin,
		FullName: sql.NullString{
			String: r.FullName,
			Valid:  r.FullName != "",
		},
		Email: sql.NullString{
			String: r.Email,
			Valid:  r.Email != "",
		},
		MustChangePassword: r.MustChangePassword,
		Active:             true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

// ToResponse converts a User model to a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:                 u.ID,
		Username:           u.Username,
		Admin:              u.Admin,
		FullName:           u.FullName.String,
		Email:              u.Email.String,
		MustChangePassword: u.MustChangePassword,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		Active:             u.Active,
	}
}
