package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/models"
)

// CreatePVZ inserts a new PVZ into the database.
func CreatePVZ(id, city, registrationDate string) error {
	query := `
    INSERT INTO pvzs (id, registration_date, city)
    VALUES (?, ?, ?)`
	_, err := DB.Exec(query, id, registrationDate, city)
	if err != nil {
		return fmt.Errorf("failed to create PVZ: %v", err)
	}
	return nil
}

// GetPVZByID retrieves a PVZ by its ID and returns it as a models.PVZ object.
func GetPVZByID(id string) (*models.PVZ, error) {
	query := `
    SELECT id, registration_date, city
    FROM pvzs
    WHERE id = ?`

	row := DB.QueryRow(query, id)

	var pvzID, registrationDate, city string
	err := row.Scan(&pvzID, &registrationDate, &city)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("PVZ not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get PVZ: %v", err)
	}

	// Parse the registration date
	parsedDate, err := time.Parse(time.RFC3339, registrationDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse registration date: %v", err)
	}

	return &models.PVZ{
		ID:               pvzID,
		RegistrationDate: parsedDate,
		City:             city,
	}, nil
}

// GetPVZsFiltered retrieves a filtered and paginated list of PVZs.
func GetPVZsFiltered(startDate, endDate string, page, limit int) ([]models.PVZ, error) {
	query := `
    SELECT id, registration_date, city
    FROM pvzs`

	var conditions []string
	var args []interface{}

	// Add filtering conditions
	if startDate != "" {
		conditions = append(conditions, "registration_date >= ?")
		args = append(args, startDate)
	}
	if endDate != "" {
		conditions = append(conditions, "registration_date <= ?")
		args = append(args, endDate)
	}

	// Join conditions with AND
	if len(conditions) > 0 {
		query += " WHERE " + joinConditions(conditions)
	}

	// Add pagination
	offset := (page - 1) * limit
	query += " ORDER BY registration_date DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// Execute the query
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get PVZs: %v", err)
	}
	defer rows.Close()

	var pvzs []models.PVZ
	for rows.Next() {
		var id, registrationDate, city string
		err := rows.Scan(&id, &registrationDate, &city)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PVZ: %v", err)
		}

		// Parse the registration date
		parsedDate, err := time.Parse(time.RFC3339, registrationDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse registration date: %v", err)
		}

		pvzs = append(pvzs, models.PVZ{
			ID:               id,
			RegistrationDate: parsedDate,
			City:             city,
		})
	}

	return pvzs, nil
}

// Helper function to join conditions with a separator
func joinConditions(conditions []string) string {
	return "(" + joinStrings(conditions, " AND ") + ")"
}

// Helper function to join strings with a separator
func joinStrings(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	for _, item := range items[1:] {
		result += sep + item
	}
	return result
}
