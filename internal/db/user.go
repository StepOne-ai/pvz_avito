package db

import (
	"database/sql"
	"fmt"
)

func CreateUser(id, email, password, role string) error {
	query := `INSERT INTO users (id, email, password, role) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, id, email, password, role)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func GetUserByEmail(email string) (map[string]string, error) {
	query := `SELECT id, email, role FROM users WHERE email = ?`
	row := DB.QueryRow(query, email)

	var id, userEmail, role string
	err := row.Scan(&id, &userEmail, &role)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return map[string]string{
		"id":    id,
		"email": userEmail,
		"role":  role,
	}, nil
}

func HashPassword(password string) string {
	return password
}
