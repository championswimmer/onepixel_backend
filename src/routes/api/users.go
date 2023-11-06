package api

import (
	"onepixel_backend/src/auth"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"

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

	// Parse incoming JSON request
	if err := ctx.BodyParser(u); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Attempt to create a new user
	newUser, err := usersController.Create(u.Email, u.Password)
	if err != nil {
		// Check if the email is already registered
		if err.Error() == "email already registered" {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
		}
		// Return 500 for all other errors
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Successfully created a new user
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registration successful!",
		"user": fiber.Map{
			"id":    newUser.ID,
			"email": newUser.Email,
		},
	})
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
