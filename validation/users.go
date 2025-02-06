// validation/users.go

package validation

import (
	"database/sql"

	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/models"
)

// UserValidator handles validation logic for users
type UserValidator struct {
	db *sql.DB
}

// NewUserValidator creates a new UserValidator instance
func NewUserValidator(db *sql.DB) *UserValidator {
	return &UserValidator{db: db}
}

// ValidateCreate performs all validations for user creation
func (v *UserValidator) ValidateCreate(user *models.UserCreateRequest) models.ValidationErrors {
	// First do basic validation
	errors := user.Validate()

	// Then check uniqueness
	if !errors.HasErrors() {
		errors = v.validateUniqueness(user.Username, user.Email)
	}

	return errors
}

func (v *UserValidator) validateUniqueness(username, email string) models.ValidationErrors {
	errors := make(models.ValidationErrors)

	// Check username uniqueness
	exists, err := db.UserExistsByUsername(v.db, username)
	if err != nil {
		errors.Add("username", "Failed to validate username uniqueness")
	} else if exists {
		errors.Add("username", "Username is already taken")
	}

	// Check email uniqueness
	if email != "" {
		exists, err := db.UserExistsByEmail(v.db, email)
		if err != nil {
			errors.Add("email", "Failed to validate email uniqueness")
		} else if exists {
			errors.Add("email", "Email is already registered")
		}
	}

	return errors
}

// ValidateUpdate performs all validations for user updates
func (v *UserValidator) ValidateUpdate(currentUser *models.User, update *models.UserUpdateRequest) models.ValidationErrors {
	errors := make(models.ValidationErrors)

	// Only validate username uniqueness if it's being changed
	if update.Username != nil && *update.Username != currentUser.Username {
		exists, err := db.UserExistsByUsernameExcept(v.db, *update.Username, currentUser.ID)
		if err != nil {
			errors.Add("username", "Failed to validate username uniqueness")
		} else if exists {
			errors.Add("username", "Username is already taken")
		}
	}

	// Only validate email uniqueness if it's being changed
	if update.Email != nil && (!currentUser.Email.Valid || *update.Email != currentUser.Email.String) {
		exists, err := db.UserExistsByEmailExcept(v.db, *update.Email, currentUser.ID)
		if err != nil {
			errors.Add("email", "Failed to validate email uniqueness")
		} else if exists {
			errors.Add("email", "Email is already registered")
		}
	}

	return errors
}
