package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"onepixel_backend/src/models"
)

// UsersRoute /api/v1/users
func UsersRoute(router fiber.Router) {
	router.Get("/", getAllUsers)
	router.Post("/", registerUser)
	router.Post("/login", loginUser)
}

func registerUser(ctx *fiber.Ctx) error {
	var u = new(models.CreateUserRequest)
	lo.Must0(ctx.BodyParser(u))

	return ctx.SendString("RegisterUser")
}

func loginUser(ctx *fiber.Ctx) error {
	return ctx.SendString("LoginUser")
}

func getAllUsers(ctx *fiber.Ctx) error {
	return ctx.SendString("GetAllUsers")
}
