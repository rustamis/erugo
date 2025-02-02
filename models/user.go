package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                 int            `json:"id"`
	Username           string         `json:"username"`
	Password           string         `json:"password_hash"`
	Admin              bool           `json:"admin"`
	FullName           sql.NullString `json:"full_name"`
	Email              sql.NullString `json:"email"`
	MustChangePassword bool           `json:"must_change_password"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Active             bool           `json:"active"`
}

type UserResponse struct {
	ID                 int       `json:"id"`
	Username           string    `json:"username"`
	Admin              bool      `json:"admin"`
	FullName           string    `json:"full_name"`
	Email              string    `json:"email"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Active             bool      `json:"active"`
}

type UserRequest struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	Admin              bool   `json:"admin"`
	FullName           string `json:"full_name"`
	Email              string `json:"email"`
	MustChangePassword bool   `json:"must_change_password"`
}

// Convert User to UserResponse
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

func (u *UserRequest) ToUser() User {
	return User{
		Username: u.Username,
		Password: u.Password,
		Admin:    u.Admin,
		FullName: sql.NullString{
			String: u.FullName,
			Valid:  u.FullName != "",
		},
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		MustChangePassword: u.MustChangePassword,
		Active:             true,
	}
}
