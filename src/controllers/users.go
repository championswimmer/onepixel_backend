package controllers

import (
	"onepixel_backend/src/db"

	"gorm.io/gorm"
)

// type UsersController
type UsersController struct {
	// db
	db *gorm.DB
}

// Create new UsersController
func NewUsersController(db *gorm.DB) *UsersController {
	return &UsersController{
		db: db,
	}
}

// create new user
func (c *UsersController) Create(email string, password string) error {
	user := &db.User{
		Email:    email,
		Password: password,
	}
	res := c.db.Create(user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
