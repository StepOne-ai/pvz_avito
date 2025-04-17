package db

import (
	"database/sql"
	"fmt"
)

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

func GetReceptionsByPVZ(pvzID string) ([]map[string]string, error) {
	query := `
    SELECT id, date_time, status FROM receptions
    WHERE pvz_id = ?
    ORDER BY date_time DESC`
	rows, err := DB.Query(query, pvzID)
	if err != nil {
		return nil, fmt.Errorf("failed to get receptions: %v", err)
	}
	defer rows.Close()

	var receptions []map[string]string
	for rows.Next() {
		var id, dateTime, status string
		err := rows.Scan(&id, &dateTime, &status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reception: %v", err)
		}
		receptions = append(receptions, map[string]string{
			"id":        id,
			"date_time": dateTime,
			"status":    status,
		})
	}
	return receptions, nil
}

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
