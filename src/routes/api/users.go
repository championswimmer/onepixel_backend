package api

import (
	"onepixel_backend/src/dtos"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

// UsersRoute /api/v1/users
func UsersRoute(router fiber.Router) {
	router.Get("/", getAllUsers)
	router.Post("/", registerUser)
	router.Post("/login", loginUser)
}

func registerUser(ctx *fiber.Ctx) error {
	var u = new(dtos.CreateUserRequest)
	lo.Must0(ctx.BodyParser(u))

	return ctx.SendString("RegisterUser")
}

func loginUser(ctx *fiber.Ctx) error {
	return ctx.SendString("LoginUser")
}

func getAllUsers(ctx *fiber.Ctx) error {
	return ctx.SendString("GetAllUsers")
}
