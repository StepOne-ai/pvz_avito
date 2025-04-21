package tests

import (
	"testing"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/db"
	"github.com/stretchr/testify/assert"
)

// TestCreatePVZInDB tests the creation of a PVZ in the database.
func TestCreatePVZInDB(t *testing.T) {
	// Initialize the database (use in-memory SQLite for testing)
	err := db.InitDB(":memory:")
	assert.NoError(t, err)

	// Test data
	pvzID := "pvz-1"
	city := "Москва"
	registrationDate := time.Now().Format(time.RFC3339) // ISO 8601 format

	// Create a PVZ
	err = db.CreatePVZ(pvzID, city, registrationDate)
	assert.NoError(t, err)

	// Retrieve the PVZ from the database
	pvz, err := db.GetPVZByID(pvzID)
	assert.NoError(t, err)
	assert.NotNil(t, pvz)

	// Validate PVZ data
	assert.Equal(t, pvzID, pvz.ID)
	assert.Equal(t, city, pvz.City)
	assert.Equal(t, registrationDate, pvz.RegistrationDate)
}
