package db

import (
	"database/sql"
	"fmt"

	"github.com/StepOne-ai/pvz_avito/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(id, email, password, role string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (id, email, password, role) VALUES (?, ?, ?, ?)`
	_, err = DB.Exec(query, id, email, hashedPassword, role)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func CheckCredentials(r_email, r_password string) error {
	query := `SELECT id, email, role, password FROM users WHERE email = ?`
	row := DB.QueryRow(query, r_email)

	var id, email, role, password string
	err := row.Scan(&id, &email, &role, &password)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user not found")
	} else if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	err = CompareHashAndPassword(password, r_password)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	return nil
}

func GetUserByEmail(email string) (models.User, error) {
	query := `SELECT id, email, role, password FROM users WHERE email = ?`
	row := DB.QueryRow(query, email)

	var id, userEmail, role, password string
	err := row.Scan(&id, &userEmail, &role, &password)
	if err == sql.ErrNoRows {
		return models.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %v", err)
	}

	return models.User{
		ID:       id,
		Email:    userEmail,
		Role:     role,
		Password: password,
	}, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
