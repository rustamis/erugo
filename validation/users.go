// validation/users.go

package validation

import (
	"database/sql"

	"github.com/DeanWard/erugo/auth"
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
func (v *UserValidator) ValidateCreate(user *models.UserCreateRequest, checkPasswordConfirmation bool) models.ValidationErrors {
	// First do basic validation
	errors := user.Validate()

	// Then check uniqueness
	if !errors.HasErrors() {
		errors = v.validateUniqueness(user.Username, user.Email)
	}

	if checkPasswordConfirmation {
		if user.Password != user.PasswordConfirmation {
			errors.Add("password_confirmation", "New password and confirmation do not match")
		}
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
func (v *UserValidator) ValidateUpdate(currentUser *models.User, update *models.UserUpdateRequest, checkCurrentPassword bool, checkPasswordConfirmation bool) models.ValidationErrors {
	errors := make(models.ValidationErrors)

	// Create a temporary UserCreateRequest to reuse validation logic
	if update.Username != nil || update.Email != nil || update.Password != nil || update.FullName != nil {
		tempUser := &models.UserCreateRequest{
			// Set current values as defaults
			Username: currentUser.Username,
			Email:    currentUser.Email.String,
			FullName: currentUser.FullName.String,
		}

		// Update with new values that are being changed
		if update.Username != nil {
			tempUser.Username = *update.Username
		}
		if update.Email != nil {
			tempUser.Email = *update.Email
		}
		if update.Password != nil {
			tempUser.Password = *update.Password
		}
		if update.FullName != nil {
			tempUser.FullName = *update.FullName
		}

		// Validate only the fields that are being updated
		createErrors := tempUser.Validate()

		if update.Password != nil {

			if checkCurrentPassword {
				if !auth.CheckPassword(currentUser.PasswordHash, *update.CurrentPassword) {
					errors.Add("current_password", "Current password is incorrect")
				}
			}
			if checkPasswordConfirmation {
				if *update.Password != *update.PasswordConfirmation {
					errors.Add("password_confirmation", "New password and confirmation do not match")
				}
			}

		}

		for field, msg := range createErrors {
			// Only include errors for fields that are actually being updated
			switch field {
			case "username":
				if update.Username != nil {
					errors.Add(field, msg)
				}
			case "email":
				if update.Email != nil {
					errors.Add(field, msg)
				}
			case "password":
				if update.Password != nil {
					errors.Add(field, msg)
				}
			case "full_name":
				if update.FullName != nil {
					errors.Add(field, msg)
				}
			}
		}

	}

	// Only proceed with uniqueness checks if basic validation passed
	if !errors.HasErrors() {
		// Check username uniqueness if it's being changed
		if update.Username != nil && *update.Username != currentUser.Username {
			exists, err := db.UserExistsByUsernameExcept(v.db, *update.Username, currentUser.ID)
			if err != nil {
				errors.Add("username", "Failed to validate username uniqueness")
			} else if exists {
				errors.Add("username", "Username is already taken")
			}
		}

		// Check email uniqueness if it's being changed
		if update.Email != nil && (!currentUser.Email.Valid || *update.Email != currentUser.Email.String) {
			exists, err := db.UserExistsByEmailExcept(v.db, *update.Email, currentUser.ID)
			if err != nil {
				errors.Add("email", "Failed to validate email uniqueness")
			} else if exists {
				errors.Add("email", "Email is already registered")
			}
		}
	}

	return errors
}
