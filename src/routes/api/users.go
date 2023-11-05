package api

import (
	"errors"
	"gorm.io/gorm"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
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
	return ctx.SendString("LoginUser")
}

func getUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("GetUserInfo")
}

func updateUserInfo(ctx *fiber.Ctx) error {
	var u = new(dtos.UpdateUserRequest)
	lo.Must0(ctx.BodyParser(u))

	if c.Locals("user").id != c.Params["id"] {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You are not allowed to perform this action",
		})
	}

	savedUser, err := usersController.Update(u.Password)
	if err != nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.UserResponseFromUser(savedUser))
}
