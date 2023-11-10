package api

import (
	"errors"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

var usersController *controllers.UsersController

// UsersRoute /api/v1/users
func UsersRoute(db *gorm.DB) func(router fiber.Router) {
	usersController = controllers.NewUsersController(db)
	return func(router fiber.Router) {
		router.Post("/", registerUser)
		router.Post("/login", loginUser)
		router.Get("/:id", auth.MandatoryAuthMiddleware, getUserInfo)
		router.Patch("/:id", auth.MandatoryAuthMiddleware, updateUserInfo)
	}
}

func registerUser(ctx *fiber.Ctx) error {
	var u = new(dtos.CreateUserRequest)
	if err := ctx.BodyParser(u); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.GetErrorResponse(fiber.StatusBadRequest, "The request body is not valid"))
	}

	savedUser, err := usersController.Create(u.Email, u.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusConflict).JSON(dtos.GetErrorResponse(fiber.StatusConflict, "User with this email already exists"))
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.UserResponseFromUser(savedUser))
}

func loginUser(ctx *fiber.Ctx) error {
	return ctx.SendString("LoginUser")
}

func getUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("GetUserInfo")
}

func updateUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("UpdateUserInfo")
}
