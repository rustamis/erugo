package validation

import (
	"database/sql"

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
	}

	if value == "" {
		errors.Add("value", "Value is required")
	}

	if group == "" {
		errors.Add("group", "Setting group is required")
	}

	// Add any specific validation rules for different setting types
	// For example, if you have boolean settings, numeric settings, etc.

	return errors
}

func (v *SettingValidator) ValidateGroup(group string) models.ValidationErrors {
	errors := make(models.ValidationErrors)

	if group == "" {
		errors.Add("group", "Setting group is required")
	}

	return errors
}
