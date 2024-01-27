package controllers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"onepixel_backend/src/db"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/security"
	"onepixel_backend/src/utils/applogger"

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

func CreateUsersController() *UsersController {
	appDb := lo.Must(db.GetAppDB())
	return &UsersController{
		db: appDb,
	}
}

func (c *UsersController) InitDefaultUser() {
	defaultUser := &models.User{
		Email:    "admin@onepixel.link",
		Password: security.HashPassword(uuid.New().String()),
		ID:       0,
		Verified: true,
	}

	// this doesn't really work, passing ID=0 to gorm is borked
	// the code here is just for documentation purposes
	res := c.db.Save([]models.User{*(defaultUser)})
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			applogger.Warn("Default user already exists")
			return
		}
		applogger.Error("Failed to create default user")
		applogger.Panic(res.Error)
	} else {
		applogger.Info("Default user created")
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
