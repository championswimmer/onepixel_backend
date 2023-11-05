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

// FindUserByEmail find user by email
func (c *UsersController) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{
		Email: email,
	}
	res := c.db.Where(user).First(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
