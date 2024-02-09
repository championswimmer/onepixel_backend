package controllers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"onepixel_backend/src/utils/applogger"
	_ "onepixel_backend/tests/providers"
	"testing"
)

var userController = CreateUsersController()

func TestUsersController_Create(t *testing.T) {
	user, token, err := userController.Create("user976@test.com", "123456")
	assert.NotNil(t, user)
	assert.NotNil(t, token)
	assert.Nil(t, err)
}

func TestUsersController_CreateDuplicateFail(t *testing.T) {
	user, token, err := userController.Create("user134534@test.com", "123456")
	assert.Nil(t, err)
	user, token, err = userController.Create("user134534@test.com", "123456")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, gorm.ErrDuplicatedKey))
}

func TestUsersController_FindUserByEmail(t *testing.T) {
	user1, _, err := userController.Create("user103439@test.com", "123456")
	assert.Nil(t, err)
	user2, err := userController.FindUserByEmail("user103439@test.com")
	assert.NotNil(t, user2)
	applogger.Info("userID", user1.ID, user2.ID)
	assert.EqualValues(t, user1.ID, user2.ID)
}

func TestCreateUserAndVerifyLogin(t *testing.T) {
	user1, token1, err := userController.Create("user139573@test.com", "123456")
	assert.Nil(t, err)
	user2, err := userController.VerifyEmailAndPassword("user139573@test.com", "123456")
	assert.NotNil(t, user2)
	assert.NotNil(t, token1)
	applogger.Info("userID", user1.ID, user2.ID)
	assert.EqualValues(t, user1.ID, user2.ID)
}
