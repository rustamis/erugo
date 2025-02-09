package validation

import (
	"database/sql"

	"github.com/DeanWard/erugo/db"
	"github.com/DeanWard/erugo/models"
)

type SettingValidator struct {
	db *sql.DB
}

func NewSettingValidator(db *sql.DB) *SettingValidator {
	return &SettingValidator{db: db}
}

func (v *SettingValidator) ValidateSet(id string, value string, group string) models.ValidationErrors {
	errors := make(models.ValidationErrors)

	if id == "" {
		errors.Add("id", "Setting ID is required")
		return errors
	}

	if value == "" {
		errors.Add("value", "Value is required")
		return errors
	}

	// Check if the setting exists
	existing, err := db.SettingById(v.db, id)
	if err != nil {
		// For new settings, we need a group
		if group == "" {
			errors.Add("group", "Setting group is required for new settings")
		}
	} else {
		// For existing settings, we ignore the group parameter as it will use the existing group
		if group != "" && group != existing.SettingGroup {
			errors.Add("group", "Cannot change setting group for existing settings")
		}
	}

	return errors
}

func (v *SettingValidator) ValidateGroup(group string) models.ValidationErrors {
	errors := make(models.ValidationErrors)

	if group == "" {
		errors.Add("group", "Setting group is required")
	}

	return errors
}
