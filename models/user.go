package models

import (
	"database/sql"
	"regexp"
	"time"
)

// ValidationErrors holds multiple validation errors by field
type ValidationErrors map[string]string

// HasErrors returns true if there are any validation errors
func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// Add adds a validation error for a field
func (v ValidationErrors) Add(field, message string) {
	v[field] = message
}

var (
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
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	Admin                bool   `json:"admin"`
	FullName             string `json:"full_name"`
	Email                string `json:"email"`
	MustChangePassword   bool   `json:"must_change_password"`
}

// UserUpdateRequest represents the expected input for user updates
type UserUpdateRequest struct {
	Username             *string `json:"username,omitempty"`
	CurrentPassword      *string `json:"current_password,omitempty"`
	Password             *string `json:"password,omitempty"`
	PasswordConfirmation *string `json:"password_confirmation,omitempty"`
	Admin                *bool   `json:"admin,omitempty"`
	FullName             *string `json:"full_name,omitempty"`
	Email                *string `json:"email,omitempty"`
	Active               *bool   `json:"active,omitempty"`
	MustChangePassword   *bool   `json:"must_change_password,omitempty"`
}

// Validate performs basic validation on the create request
func (r *UserCreateRequest) Validate() ValidationErrors {
	errors := make(ValidationErrors)

	// Username validation
	if r.Username == "" {
		errors.Add("username", "Username is required")
	} else if len(r.Username) < 3 {
		errors.Add("username", "Username must be at least 3 characters long")
	}

	// Password validation
	if r.Password == "" {
		errors.Add("password", "Password is required")
	} else if len(r.Password) < 8 {
		errors.Add("password", "Password must be at least 8 characters long")
	}

	// Full name validation
	if r.FullName == "" {
		errors.Add("full_name", "Full name is required")
	}

	// Email validation
	if r.Email == "" {
		errors.Add("email", "Email is required")
	} else if !emailRegex.MatchString(r.Email) {
		errors.Add("email", "Email is invalid")
	}

	return errors
}

// ToUser converts a creation request to a User model
func (r *UserCreateRequest) ToUser() User {
	return User{
		Username:     r.Username,
		PasswordHash: r.Password, // Will be hashed by the db package
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
