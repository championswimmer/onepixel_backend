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

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute)
	apiV1.Route("/urls", api.UrlsRoute)

	log.Fatal(app.Listen(":3000"))
}
