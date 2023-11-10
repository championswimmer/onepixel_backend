package api

import (
	"errors"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

// UsersRoute /api/v1/users
func UsersRoute(router fiber.Router, usersController *controllers.UsersController) {
	router.Get("/", getAllUsers)
	router.Post("/", func(c *fiber.Ctx) error {
		return registerUser(c, usersController)
	})
	router.Post("/login", loginUser)
	router.Get("/:id", auth.MandatoryAuthMiddleware, getUserInfo)
	router.Patch("/:id", auth.MandatoryAuthMiddleware, updateUserInfo)
}

func registerUser(ctx *fiber.Ctx, usersController *controllers.UsersController) error {
	var u = new(dtos.CreateUserRequest)
	if err := ctx.BodyParser(u); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.CreateErrorResponse(
			fiber.StatusBadRequest,
			"The request body is not valid",
		))
	}

	if u.Email == "" || u.Password == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dtos.CreateErrorResponse(
			fiber.StatusUnprocessableEntity,
			"email and password are required to create user",
		))
	}

	savedUser, err := usersController.Create(u.Email, u.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusConflict).JSON(dtos.CreateErrorResponse(fiber.StatusConflict, "User with this email already exists"))
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUserResponseFromUser(savedUser))
}

func loginUser(ctx *fiber.Ctx) error {
	return ctx.SendString("LoginUser")
}

func getAllUsers(ctx *fiber.Ctx) error {
	return ctx.SendString("GetAllUsers")
}

func getUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("GetUserInfo")
}

func updateUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("UpdateUserInfo")
}
