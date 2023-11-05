package api

import (
  "errors"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/middleware"

	"gorm.io/gorm"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

var usersController *controllers.UsersController

// UsersRoute /api/v1/users
/*
	Header structure for withAuthRouter routes

	Authorization: Bearer <jwt_token>
*/
func UsersRoute(db *gorm.DB) func(router fiber.Router) {
	usersController = controllers.NewUsersController(db)
	return func(withoutAuthRouter fiber.Router) {
		withoutAuthRouter.Post("/", registerUser)
		withoutAuthRouter.Post("/login", loginUser)
		withoutAuthRouter.Get("/:id", getUserInfo)

		withAuthRouter := withoutAuthRouter.Group("/auth", middleware.RequireJwtAuth())
		withAuthRouter.Patch("/:id", updateUserInfo)
	}
}

func registerUser(ctx *fiber.Ctx) error {
	var u = new(dtos.CreateUserRequest)
	// TODO: handle error and show Bad-Request to client
	lo.Must0(ctx.BodyParser(u))

	savedUser, err := usersController.Create(u.Email, u.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			// TODO: make a Error response DTO
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "User with this email already exists",
			})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.UserResponseFromUser(savedUser))
}

func loginUser(ctx *fiber.Ctx) error {
	var u = new(dtos.CreateUserRequest)
	lo.Must0(ctx.BodyParser(u))
	savedUser := lo.Must(usersController.Get(u.Email, u.Password))

	// Throws Unauthorized error
	if u.Password != savedUser.Password {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseFromServer("Incorrect credentials"))
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"email": savedUser.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("some_random_secret_string"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseFromServer("Internal server error"))
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.LoginResponseFromUser(t))
}

func getUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("GetUserInfo")
}

func updateUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("UpdateUserInfo")
}
