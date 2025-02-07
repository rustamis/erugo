package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/DeanWard/erugo/models"
)

// Common database errors
var (
	ErrSettingNotFound = fmt.Errorf("setting not found")
	ErrSettingInvalid  = fmt.Errorf("setting invalid")
)

func SettingsByGroup(db *sql.DB, group string) ([]models.Setting, error) {
	query := `
        SELECT id, value, previous_value, setting_group 
        FROM settings 
        WHERE setting_group = ?
        ORDER BY id`

	rows, err := db.Query(query, group)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}
	defer rows.Close()

	var settings []models.Setting
	for rows.Next() {
		var setting models.Setting
		err := rows.Scan(
			&setting.Id,
			&setting.Value,
			&setting.PreviousValue,
			&setting.SettingGroup,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
		}
		settings = append(settings, setting)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	return settings, nil
}

func SettingById(db *sql.DB, id string) (*models.Setting, error) {
	query := `
        SELECT id, value, previous_value, setting_group 
        FROM settings 
        WHERE id = ?`

	var setting models.Setting
	err := db.QueryRow(query, id).Scan(
		&setting.Id,
		&setting.Value,
		&setting.PreviousValue,
		&setting.SettingGroup,
	)

	if err == sql.ErrNoRows {
		return nil, ErrSettingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	return &setting, nil
}

func SettingSetById(db *sql.DB, id string, newValue string, group string) error {
	// First try to get the current value
	current, err := SettingById(db, id)

	if err != nil {
		// If setting doesn't exist, create it
		if errors.Is(err, ErrSettingNotFound) {
			query := `
                INSERT INTO settings (id, value, previous_value, setting_group, created_at, updated_at)
                VALUES (?, ?, '', ?, ?, ?)`

			now := time.Now()
			result, err := db.Exec(query, id, newValue, group, now, now)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrDBOperation, err)
			}

			rows, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("%w: %v", ErrDBOperation, err)
			}
			if rows == 0 {
				return fmt.Errorf("%w: failed to create setting", ErrDBOperation)
			}

			return nil
		}
		return err
	}

	// Setting exists, update it
	query := `
        UPDATE settings 
        SET value = ?, 
            previous_value = ?,
            updated_at = ?
        WHERE id = ?`

	result, err := db.Exec(query, newValue, current.Value, time.Now(), id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBOperation, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDBOperation, err)
	}
	if rows == 0 {
		return ErrSettingNotFound
	}

	return nil
}
