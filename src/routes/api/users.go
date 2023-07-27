package api

import "github.com/gofiber/fiber/v2"

func UsersRoute(router fiber.Router) {
	router.Get("/", getAllUsers)
	router.Post("/register", registerUser)
	router.Post("/login", loginUser)
}

func registerUser(ctx *fiber.Ctx) error {
	return ctx.SendString("RegisterUser")
}

func loginUser(ctx *fiber.Ctx) error {
	return ctx.SendString("LoginUser")
}

func getAllUsers(ctx *fiber.Ctx) error {
	return ctx.SendString("GetAllUsers")
}
