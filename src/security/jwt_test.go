package security

import (
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils/applogger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJwt_CreateToken(t *testing.T) {
	testUser := &models.User{
		ID: 12,
	}

	jwt := CreateJWTFromUser(testUser)
	applogger.Info("CreateToken", "jwt", jwt)
	assert.NotNil(t, jwt)
}

func TestJwt_ParseToken(t *testing.T) {
	testUser := &models.User{
		ID: 12,
	}

	jwt := CreateJWTFromUser(testUser)
	applogger.Info("CreateToken", "jwt", jwt)
	assert.NotNil(t, jwt)

	user, err := ValidateJWT(jwt)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.ID, user.ID)
}
