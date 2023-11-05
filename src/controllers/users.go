package controllers

import (
	"errors"
	"onepixel_backend/src/models"

	"github.com/samber/lo"
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
func (c *UsersController) Create(email string, password string) (*models.User, error) {
	// Check if email is already registered
	existingUser, err := c.FindUserByEmail(email)
	if err == nil && existingUser.ID != 0 {
		// User exists and ID is populated, hence email is already registered
		return nil, errors.New("email already registered")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// There was an actual error in looking up the user
		return nil, err
	}

	// Hash the password
	hashedPassword := lo.Must(bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost))

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save the user to the database
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
