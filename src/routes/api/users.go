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
		// Public Routes
		router.Post("/", registerUser)
		router.Post("/login", loginUser)

		// Private Routes
		router.Get("/:id", auth.MandatoryAuthMiddleware, getUserInfo)
		router.Patch("/:id", auth.MandatoryAuthMiddleware, updateUserInfo)
	}
}

func registerUser(ctx *fiber.Ctx) error {
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
	var u = new(dtos.CreateUserRequest)
	lo.Must0(ctx.BodyParser(u))
	savedUser := lo.Must(usersController.FindUserByEmail(u.Email))

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
