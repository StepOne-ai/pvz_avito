package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StepOne-ai/pvz_avito/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	r.POST("/protected", middleware.JWTMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	// Generate a valid JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   "test-user-id",
		"role": "employee",
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)

	// Create a test HTTP request with the token
	req, _ := http.NewRequest(http.MethodPost, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoleMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	r := gin.Default()
	r.POST("/moderator-only", middleware.JWTMiddleware(), middleware.RoleMiddleware("moderator"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	// Generate a valid JWT token for an employee
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   "test-user-id",
		"role": "employee",
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)

	// Create a test HTTP request with the token
	req, _ := http.NewRequest(http.MethodPost, "/moderator-only", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response (access should be denied)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
