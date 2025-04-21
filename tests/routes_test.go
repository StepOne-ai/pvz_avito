package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/StepOne-ai/pvz_avito/internal/db"
	"github.com/StepOne-ai/pvz_avito/internal/middleware"
	"github.com/StepOne-ai/pvz_avito/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestRegisterEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	routes.SetupRoutes(r)

	reqBody := `{
        "email": "test@example.com",
        "password": "securepassword",
        "role": "employee"
    }`

	// Create a test HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	var responseBody map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", responseBody["email"])
	assert.Equal(t, "employee", responseBody["role"])
}

func TestLoginEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Register a user
	db.InitDB(":memory:")
	err := db.CreateUser("test-user-id", "test@example.com", "securepassword", "employee")
	assert.NoError(t, err)

	// Test valid login
	validReqBody := `{
        "email": "test@example.com",
        "password": "securepassword"
    }`
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(validReqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test invalid login
	invalidReqBody := `{
        "email": "test@example.com",
        "password": "wrongpassword"
    }`
	req, _ = http.NewRequest(http.MethodPost, "/login", strings.NewReader(invalidReqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreatePVZ(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Generate a valid JWT token for a moderator
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   "moderator-id",
		"role": "moderator",
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)

	// Test request body
	reqBody := `{
        "city": "Москва"
    }`

	// Create a test HTTP request with the token
	req, _ := http.NewRequest(http.MethodPost, "/pvz", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Validate the response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Москва", responseBody["city"])
}

func TestGetPVZList(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Generate a valid JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   "user-id",
		"role": "moderator",
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)

	// Add test PVZs to the database
	db.InitDB(":memory:")
	db.CreatePVZ("pvz-1", "Москва", "2023-01-01T00:00:00Z")
	db.CreatePVZ("pvz-2", "Санкт-Петербург", "2023-02-01T00:00:00Z")

	// Create a test HTTP request
	req, _ := http.NewRequest(http.MethodGet, "/pvz?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Validate the response body
	var responseBody []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Len(t, responseBody, 2)
}

func TestCloseLastReception(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Generate a valid JWT token for a moderator
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   "moderator-id",
		"role": "moderator",
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)

	// Add a test PVZ and reception to the database
	db.InitDB(":memory:")
	db.CreatePVZ("pvz-1", "Москва", "2023-01-01T00:00:00Z")
	db.CreateReception("reception-1", "2023-01-02T00:00:00Z", "pvz-1", "in_progress")

	// Create a test HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/pvz/pvz-1/close_last_reception", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Validate the response body
	var responseBody map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "close", responseBody["status"])
}

func TestDeleteLastProduct(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Generate a valid JWT token for an employee
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   "employee-id",
		"role": "employee",
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)

	// Add a test PVZ, reception, and product to the database
	db.InitDB(":memory:")
	db.CreatePVZ("pvz-1", "Москва", "2023-01-01T00:00:00Z")
	db.CreateReception("reception-1", "2023-01-02T00:00:00Z", "pvz-1", "in_progress")
	db.CreateProduct("product-1", "2023-01-03T00:00:00Z", "электроника", "reception-1")

	// Create a test HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/pvz/pvz-1/delete_last_product", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Validate the response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Товар удален", responseBody["message"])
}

func TestCreateReception(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Generate a valid JWT token for an employee
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   "employee-id",
		"role": "employee",
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)

	// Add a test PVZ to the database
	db.InitDB(":memory:")
	db.CreatePVZ("pvz-1", "Москва", "2023-01-01T00:00:00Z")

	// Test request body
	reqBody := `{
        "pvzId": "pvz-1"
    }`

	// Create a test HTTP request with the token
	req, _ := http.NewRequest(http.MethodPost, "/receptions", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Validate the response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "pvz-1", responseBody["pvzId"])
}

func TestAddProduct(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	routes.SetupRoutes(r)

	// Generate a valid JWT token for an employee
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   "employee-id",
		"role": "employee",
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)

	// Add a test PVZ and reception to the database
	db.InitDB(":memory:")
	db.CreatePVZ("pvz-1", "Москва", "2023-01-01T00:00:00Z")
	db.CreateReception("reception-1", "2023-01-02T00:00:00Z", "pvz-1", "in_progress")

	// Test request body
	reqBody := `{
        "type": "электроника",
        "pvzId": "pvz-1"
    }`

	// Create a test HTTP request with the token
	req, _ := http.NewRequest(http.MethodPost, "/products", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Validate the response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "электроника", responseBody["type"])
}
