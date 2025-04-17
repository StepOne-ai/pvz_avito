package db

import (
	"database/sql"
	"fmt"
)

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

func GetPVZByID(id string) (map[string]string, error) {
	query := `
    SELECT id, registration_date, city FROM pvzs WHERE id = ?`
	row := DB.QueryRow(query, id)

	var pvzID, registrationDate, city string
	err := row.Scan(&pvzID, &registrationDate, &city)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("PVZ not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get PVZ: %v", err)
	}

	return map[string]string{
		"id":                pvzID,
		"registration_date": registrationDate,
		"city":              city,
	}, nil
}

func GetPVZsFiltered(startDate, endDate string, page, limit int) ([]map[string]interface{}, error) {
	// Base query
	query := `
    SELECT id, registration_date, city FROM pvzs`

	var conditions []string
	var args []any

	if startDate != "" {
		conditions = append(conditions, "registration_date >= ?")
		args = append(args, startDate)
	}
	if endDate != "" {
		conditions = append(conditions, "registration_date <= ?")
		args = append(args, endDate)
	}

	if len(conditions) > 0 {
		query += " WHERE " + joinConditions(conditions)
	}

	offset := (page - 1) * limit
	query += " ORDER BY registration_date DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get PVZs: %v", err)
	}
	defer rows.Close()

	var pvzs []map[string]any
	for rows.Next() {
		var id, registrationDate, city string
		err := rows.Scan(&id, &registrationDate, &city)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PVZ: %v", err)
		}

		pvzs = append(pvzs, map[string]any{
			"id":                id,
			"registration_date": registrationDate,
			"city":              city,
			"receptions":        []any{},
		})
	}

	return pvzs, nil
}

func joinConditions(conditions []string) string {
	return "(" + joinStrings(conditions, " AND ") + ")"
}

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
