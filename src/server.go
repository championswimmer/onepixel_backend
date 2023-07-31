package main

import (
	"log"
	"onepixel_backend/src/db"
	"onepixel_backend/src/routes/api"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute)
	apiV1.Route("/urls", api.UrlsRoute)

	// Initialize the database
	lo.Must(db.InitDB())

	log.Fatal(app.Listen(":3000"))
}
