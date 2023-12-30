package controllers

import (
	"onepixel_backend/src/models"
	"onepixel_backend/src/security"

	"gorm.io/gorm"
)

type AuthError struct {
	reason string
}

func (e *AuthError) Error() string {
	return e.reason
}

// UserController errors
var (
	PasswordInvalidLoginError = &AuthError{"Invalid password"}
	EmailInvalidLoginError    = &AuthError{"Invalid email"}
)

type UsersController struct {
	// db
	db *gorm.DB
}

func CreateUsersController(db *gorm.DB) *UsersController {
	return &UsersController{
		db: db,
	}
}

// Create new user
func (c *UsersController) Create(email string, password string) (user *models.User, token string, err error) {
	user = &models.User{
		Email:    email,
		Password: security.HashPassword(password),
	}
	res := c.db.Create(user)
	if res.Error != nil {
		return nil, "", res.Error
	}
	token = security.CreateJWTFromUser(user)
	return
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

func (c *UsersController) VerifyEmailAndPassword(email string, password string) (*models.User, error) {
	user, err := c.FindUserByEmail(email)
	if err != nil {
		return nil, EmailInvalidLoginError
	}
	if !security.CheckPasswordHash(password, user.Password) {
		return nil, PasswordInvalidLoginError
	}
	return user, nil
}
