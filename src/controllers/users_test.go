package controllers

import (
	"errors"
	"log"
	"onepixel_backend/src/db"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var userController = NewUsersController(lo.Must(db.InitDBTest()))

func TestUsersController_Create(t *testing.T) {
	user, err := userController.Create("user976@test.com", "123456")
	assert.NotNil(t, user)
	assert.Nil(t, err)
}

func TestUsersController_CreateDuplicateFail(t *testing.T) {
	user, err := userController.Create("user134534@test.com", "123456")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	user, err = userController.Create("user134534@test.com", "123456")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, gorm.ErrDuplicatedKey))
}

func TestUsersController_FindUserByEmail(t *testing.T) {
	user1, err := userController.Create("user103439@test.com", "123456")
	assert.Nil(t, err)
	user2, err := userController.FindUserByEmail("user103439@test.com")
	assert.NotNil(t, user2)
	assert.Nil(t, err)
	log.Println("userID", user1.ID, user2.ID)
	assert.EqualValues(t, user1.ID, user2.ID)
}

func TestUsersController_FindUserByEmailNonExistent(t *testing.T) {
    user, err := userController.FindUserByEmail("nonexistent@test.com")
    assert.Nil(t, user)
    assert.NotNil(t, err)
    assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}


func TestUsersController_CreatePasswordHashing(t *testing.T) {
    plainPassword := "123456"
    user, err := userController.Create("userhash@test.com", plainPassword)
    assert.Nil(t, err)
    assert.NotEqual(t, plainPassword, user.Password)
}


func TestUsersController_CreateInvalidEmail(t *testing.T) {
    user, err := userController.Create("", "123456") // Testing with empty email
    assert.Nil(t, user)
    assert.NotNil(t, err)
}


func TestUsersController_CreateUserAttributeIntegrity(t *testing.T) {
    email := "integrity@test.com"
    password := "123456"
    user, err := userController.Create(email, password)
    assert.Nil(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, email, user.Email)
}
