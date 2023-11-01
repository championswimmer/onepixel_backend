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
func (c *UsersController) Create(email string, password string) error {
	user := &models.User{
		Email:    email,
		Password: password,
	}
	res := c.db.Create(user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
