package api

import (
	"errors"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

var usersController *controllers.UsersController

// UsersRoute /api/v1/users
func UsersRoute(db *gorm.DB) func(router fiber.Router) {
	usersController = controllers.NewUsersController(db)
	return func(router fiber.Router) {
		// Public Routes
		router.Post("/", registerUser)
		router.Post("/login", loginUser)
		router.Get("/:id", getUserInfo)

		// Private Routes
		router.Patch("/:id", auth.MandatoryAuthMiddleware, updateUserInfo)
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
	savedUser := lo.Must(usersController.Get(u.Email))

	if u.Password != savedUser.Password {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseFromServer("Incorrect credentials"))
	}
	token := auth.CreateJWTFromUser(savedUser)
	return ctx.Status(fiber.StatusOK).JSON(dtos.LoginResponseFromUser(token))
}

func getUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("GetUserInfo")
}

func updateUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("UpdateUserInfo")
}
