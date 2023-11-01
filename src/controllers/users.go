package controllers

import (
	"errors"
	"onepixel_backend/src/models"

	"golang.org/x/crypto/bcrypt"
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
	// Check if email is already registered
	existingUser := &models.User{}
	result := c.db.Where("email = ?", email).First(existingUser)
	if result.Error == nil {
		return errors.New("email already registered")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save the user to the database
	res := c.db.Create(user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
