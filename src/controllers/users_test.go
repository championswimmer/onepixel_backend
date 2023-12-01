package controllers

import (
	"errors"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"onepixel_backend/src/db"
	"testing"
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
	log.Println("userID", user1.ID, user2.ID)
	assert.EqualValues(t, user1.ID, user2.ID)
}

func TestCreateUserAndVerifyLogin(t *testing.T) {
	user1, err := userController.Create("user139573@test.com", "123456")
	assert.Nil(t, err)
	user2, err := userController.VerifyEmailAndPassword("user139573@test.com", "123456")
	assert.NotNil(t, user2)
	log.Println("userID", user1.ID, user2.ID)
	assert.EqualValues(t, user1.ID, user2.ID)
}
