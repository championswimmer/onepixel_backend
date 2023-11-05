package auth

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/stretchr/testify/assert"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/models"
	"testing"
)

func TestJwt_CreateToken(t *testing.T) {
	testUser := &models.User{
		ID: 12,
	}

	jwt := auth.CreateJWTFromUser(testUser)
	log.Info("jwt: ", jwt)
	assert.NotNil(t, jwt)
}

func TestJwt_ParseToken(t *testing.T) {
	testUser := &models.User{
		ID: 12,
	}

	jwt := auth.CreateJWTFromUser(testUser)
	log.Info("jwt: ", jwt)
	assert.NotNil(t, jwt)

	userID, err := auth.ValidateJWT(jwt)
	assert.Nil(t, err)
	assert.NotNil(t, userID)
	assert.Equal(t, testUser.ID, *userID)
}
