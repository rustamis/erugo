package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Admin        bool   `json:"admin"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

// Convert User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Admin:    u.Admin,
	}
}
