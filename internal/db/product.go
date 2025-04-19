package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/models"
)

// CreateProduct inserts a new product into the database.
func CreateProduct(id, dateTime, productType, receptionId string) error {
	query := `
    INSERT INTO products (id, date_time, type, reception_id)
    VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, id, dateTime, productType, receptionId)
	if err != nil {
		return fmt.Errorf("failed to create product: %v", err)
	}
	return nil
}

// GetProductsByReception retrieves all products for a given reception ID.
func GetProductsByReception(receptionId string) ([]models.Product, error) {
	query := `
    SELECT id, date_time, type
    FROM products
    WHERE reception_id = ?
    ORDER BY date_time DESC`
	rows, err := DB.Query(query, receptionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var id, dateTime, productType string
		err := rows.Scan(&id, &dateTime, &productType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}

		// Parse the date_time field
		parsedTime, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %v", err)
		}

		products = append(products, models.Product{
			ID:          id,
			DateTime:    parsedTime,
			Type:        productType,
			ReceptionId: receptionId,
		})
	}
	return products, nil
}

// DeleteLastProduct deletes the last product added to the current reception for a given PVZ ID.
func DeleteLastProduct(pvzId string) error {
	// Find the active reception for the given PVZ ID
	receptionQuery := `
    SELECT id
    FROM receptions
    WHERE pvz_id = ? AND status = 'in_progress'
    ORDER BY date_time DESC
    LIMIT 1`
	var receptionId string
	err := DB.QueryRow(receptionQuery, pvzId).Scan(&receptionId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no active reception found for PVZ ID: %s", pvzId)
	} else if err != nil {
		return fmt.Errorf("failed to find active reception: %v", err)
	}

	// Find the last product in the active reception
	productQuery := `
    SELECT id
    FROM products
    WHERE reception_id = ?
    ORDER BY date_time DESC
    LIMIT 1`
	var productId string
	err = DB.QueryRow(productQuery, receptionId).Scan(&productId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no products found in the current reception")
	} else if err != nil {
		return fmt.Errorf("failed to find last product: %v", err)
	}

	// Delete the product
	deleteQuery := `
    DELETE FROM products
    WHERE id = ?`
	_, err = DB.Exec(deleteQuery, productId)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}
	return nil
}
