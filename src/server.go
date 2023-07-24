package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"onepixel_backend/src/routes/api"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	api.Users(app.Group("/api/users"))

	log.Fatal(app.Listen(":3000"))
}
