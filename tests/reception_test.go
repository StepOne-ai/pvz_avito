package tests

import (
	"testing"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateReceptionInDB(t *testing.T) {
	err := db.InitDB(":memory:")
	assert.NoError(t, err)

	// Add a test PVZ to the database
	pvzID := "pvz-1"
	city := "Москва"
	registrationDate := time.Now().Format(time.RFC3339)
	err = db.CreatePVZ(pvzID, city, registrationDate)
	assert.NoError(t, err)

	// Test data
	receptionID := "reception-1"
	dateTime := time.Now().Format(time.RFC3339)
	status := "in_progress"

	// Create a reception
	err = db.CreateReception(receptionID, dateTime, pvzID, status)
	assert.NoError(t, err)

	// Retrieve the reception from the database
	reception, err := db.GetReceptionByID(receptionID)
	assert.NoError(t, err)
	assert.NotNil(t, reception)

	// Validate reception data
	assert.Equal(t, receptionID, reception.ID)
	assert.Equal(t, dateTime, reception.DateTime)
	assert.Equal(t, pvzID, reception.PvzId)
	assert.Equal(t, status, reception.Status)
}
