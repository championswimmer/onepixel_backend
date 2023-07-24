package api

import "github.com/gofiber/fiber/v2"

func registerUser(ctx *fiber.Ctx) error {
	return ctx.SendString("RegisterUser")
}

func loginUser(ctx *fiber.Ctx) error {
	return ctx.SendString("LoginUser")
}

func Users(router fiber.Router) {
	router.Post("/register", registerUser)
	router.Post("/login", loginUser)
}
