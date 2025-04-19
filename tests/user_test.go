package tests

import (
	"testing"

	"github.com/StepOne-ai/pvz_avito/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetUser(t *testing.T) {
	err := db.InitDB(":memory:")
	assert.NoError(t, err)

	userID := "test-user-id"
	email := "test@example.com"
	password := "securepassword"
	role := "employee"

	err = db.CreateUser(userID, email, password, role)
	assert.NoError(t, err)

	user, err := db.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, role, user.Role)
}
