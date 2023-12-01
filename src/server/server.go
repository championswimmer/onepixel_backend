package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
	_ "onepixel_backend/src/docs"
	"onepixel_backend/src/routes/api"
)

func CreateApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute(db))
	apiV1.Route("/urls", api.UrlsRoute)

	app.Get("/docs/*", swagger.HandlerDefault)

	return app
}
