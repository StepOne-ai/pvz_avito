package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/models"
)

// CreateReception inserts a new reception into the database.
func CreateReception(id, dateTime, pvzId, status string) error {
	query := `
    INSERT INTO receptions (id, date_time, pvz_id, status)
    VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, id, dateTime, pvzId, status)
	if err != nil {
		return fmt.Errorf("failed to create reception: %v", err)
	}
	return nil
}

// GetReceptionsByPVZ retrieves all receptions for a given PVZ ID.
func GetReceptionsByPVZ(pvzID string) ([]models.Reception, error) {
	query := `
    SELECT id, date_time, status FROM receptions
    WHERE pvz_id = ?
    ORDER BY date_time DESC`
	rows, err := DB.Query(query, pvzID)
	if err != nil {
		return []models.Reception{}, fmt.Errorf("failed to get receptions: %v", err)
	}
	defer rows.Close()

	var receptions []models.Reception
	for rows.Next() {
		var id, dateTime, status string
		err := rows.Scan(&id, &dateTime, &status)
		if err != nil {
			return []models.Reception{}, fmt.Errorf("failed to scan reception: %v", err)
		}
		parsedTime, err := time.Parse("2006-01-02 15:04:05 MST", dateTime)
		if err != nil {
			return []models.Reception{}, fmt.Errorf("failed to parse time: %v", err)
		}

		receptions = append(receptions, models.Reception{
			ID:       id,
			DateTime: parsedTime,
			Status:   status,
		})
	}
	return receptions, nil
}

// CloseLastReception closes the last active reception for a given PVZ ID.
func CloseLastReception(pvzId string) error {
	query := `
    SELECT id FROM receptions
    WHERE pvz_id = ? AND status = 'in_progress'
    ORDER BY date_time DESC
    LIMIT 1`
	var receptionId string
	err := DB.QueryRow(query, pvzId).Scan(&receptionId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no active reception found for PVZ ID: %s", pvzId)
	} else if err != nil {
		return fmt.Errorf("failed to find active reception: %v", err)
	}

	updateQuery := `
    UPDATE receptions
    SET status = 'close'
    WHERE id = ?`
	_, err = DB.Exec(updateQuery, receptionId)
	if err != nil {
		return fmt.Errorf("failed to close reception: %v", err)
	}

	return nil
}

// GetReceptionByID retrieves a reception by its ID and returns it as a models.Reception object.
func GetReceptionByID(id string) (*models.Reception, error) {
	query := `
    SELECT id, date_time, pvz_id, status
    FROM receptions
    WHERE id = ?`

	row := DB.QueryRow(query, id)

	// Variables to store the result
	var receptionID, dateTime, pvzID, status string

	// Scan the row into variables
	err := row.Scan(&receptionID, &dateTime, &pvzID, &status)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("reception not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get reception: %v", err)
	}

	// Parse the date_time field
	parsedTime, err := time.Parse("2006-01-02 15:04:05 MST", dateTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	// Return the reception data as a models.Reception object
	return &models.Reception{
		ID:       receptionID,
		DateTime: parsedTime,
		PvzId:    pvzID,
		Status:   status,
	}, nil
}
