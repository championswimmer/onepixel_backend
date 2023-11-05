package controllers

import (
	"onepixel_backend/src/models"

	"gorm.io/gorm"
)

type UsersController struct {
	// db
	db *gorm.DB
}

func NewUsersController(db *gorm.DB) *UsersController {
	return &UsersController{
		db: db,
	}
}

// Get User
func (c *UsersController) Get(email string, password string) (*models.User, error) {
	user := &models.User{
		Email:    email,
		Password: password, // TODO: hash password
	}
	res := c.db.First(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

// Create new user
func (c *UsersController) Create(email string, password string) (*models.User, error) {
	user := &models.User{
		Email:    email,
		Password: password, // TODO: hash password
	}
	res := c.db.Create(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
