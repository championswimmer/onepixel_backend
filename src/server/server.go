package server

import (
	"onepixel_backend/src/routes/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func CreateApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute(db))
	apiV1.Route("/urls", api.UrlsRoute)

	return app
}
