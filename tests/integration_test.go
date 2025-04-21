package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/db"
	"github.com/StepOne-ai/pvz_avito/internal/models"
	"github.com/StepOne-ai/pvz_avito/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_PVZReceptionAndProducts tests the integration flow:
// 1. Creates a new PVZ
// 2. Adds a new reception
// 3. Adds 50 products to the reception
// 4. Closes the reception
func TestIntegration_PVZReceptionAndProducts(t *testing.T) {
	// Initialize the Gin router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Initialize the database (in-memory SQLite)
	err := db.InitDB(":memory:")
	require.NoError(t, err)

	// Step 1: Create a new PVZ
	pvzID := "pvz-123"
	city := "Москва"
	reqBodyPVZ := models.PVZ{
		ID:               pvzID,
		City:             city,
		RegistrationDate: time.Now(),
	}

	// Send POST request to create PVZ
	respPVZ, bodyPVZ := makeRequest(t, http.MethodPost, "/pvz", reqBodyPVZ, r)
	assert.Equal(t, http.StatusCreated, respPVZ.StatusCode)
	var createdPVZ models.PVZ
	err = json.Unmarshal(bodyPVZ, &createdPVZ)
	require.NoError(t, err)
	assert.Equal(t, pvzID, createdPVZ.ID)
	assert.Equal(t, city, createdPVZ.City)

	// Step 2: Add a new reception
	receptionID := "reception-456"
	reqBodyReception := map[string]string{
		"pvzId": pvzID,
	}

	// Send POST request to create reception
	respReception, bodyReception := makeRequest(t, http.MethodPost, "/receptions", reqBodyReception, r)
	assert.Equal(t, http.StatusCreated, respReception.StatusCode)
	var createdReception models.Reception
	err = json.Unmarshal(bodyReception, &createdReception)
	require.NoError(t, err)
	assert.Equal(t, receptionID, createdReception.ID)
	assert.Equal(t, pvzID, createdReception.PvzId)

	// Step 3: Add 50 products to the reception
	for i := 1; i <= 50; i++ {
		productType := fmt.Sprintf("Товар %d", i)
		reqBodyProduct := map[string]string{
			"type":  productType,
			"pvzId": pvzID,
		}

		// Send POST request to add a product
		respProduct, _ := makeRequest(t, http.MethodPost, "/products", reqBodyProduct, r)
		assert.Equal(t, http.StatusCreated, respProduct.StatusCode)
	}

	// Step 4: Close the reception
	reqBodyClose := map[string]string{
		"pvzId": pvzID,
	}

	// Send POST request to close the reception
	respClose, bodyClose := makeRequest(t, http.MethodPost, "/pvz/"+pvzID+"/close_last_reception", reqBodyClose, r)
	assert.Equal(t, http.StatusOK, respClose.StatusCode)
	var closedReception models.Reception
	err = json.Unmarshal(bodyClose, &closedReception)
	require.NoError(t, err)
	assert.Equal(t, pvzID, closedReception.PvzId)
	assert.Equal(t, "close", closedReception.Status)
}

// Helper function to make HTTP requests
func makeRequest(t *testing.T, method string, url string, body interface{}, r *gin.Engine) (*http.Response, []byte) {
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w.Result(), w.Body.Bytes()
}
