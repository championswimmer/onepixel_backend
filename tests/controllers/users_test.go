package controllers

import (
	"errors"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/db"
	"testing"
)

var userController = controllers.NewUsersController(lo.Must(db.InitDBTest()))

func TestUsersController_Create(t *testing.T) {
	user, err := userController.Create("user@test.com", "123456")
	assert.NotNil(t, user)
	assert.Nil(t, err)
}

func TestUsersController_CreateDuplicateFail(t *testing.T) {
	user, err := userController.Create("user2@test.com", "123456")
	assert.Nil(t, err)
	user, err = userController.Create("user2@test.com", "123456")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.True(t, errors.Is(err, gorm.ErrDuplicatedKey))
}

func TestUsersController_FindUserByEmail(t *testing.T) {
	user, err := userController.Create("user2@test.com", "123456")
	assert.Nil(t, err)
	userId := user.ID
	user, err = userController.FindUserByEmail("user2@test.com")
	assert.NotNil(t, user)
	assert.Equal(t, userId, user.ID)
}
