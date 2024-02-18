package controllers

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onepixel_backend/src/config"
	"onepixel_backend/src/db"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/security"
	"onepixel_backend/src/utils/applogger"
	"sync"
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

var initDefaultUserOnce sync.Once

func (c *UsersController) initDefaultUser() {
	var defaultUser = models.User{
		Email: config.AdminUserEmail,
	}
	// check if default user exists
	res := c.db.First(&defaultUser)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			applogger.Info("Default user doesn't exist")
			defaultUser.Password = security.HashPassword(uuid.New().String())
			defaultUser.Verified = true
			res = c.db.Create(&defaultUser)
		} else {
			applogger.Error("Failed to check if default user exists")
			applogger.Panic(res.Error)
		}
	}
	// update default user id to 0
	res = c.db.Model(defaultUser).
		Where("email = ?", defaultUser.Email).
		Update("id", 0)
	if res.Error != nil {
		applogger.Error("Failed to update default user id")
		applogger.Panic(res.Error)
	}

}

func CreateUsersController() *UsersController {
	appDb := db.GetAppDB()
	ctrl := &UsersController{
		db: appDb,
	}
	initDefaultUserOnce.Do(ctrl.initDefaultUser)
	return ctrl
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
