package server

import (
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/routes/api"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

func CreateApp(dbConn *gorm.DB) *fiber.App {
	app := fiber.New()

	usersController := controllers.NewUsersController(dbConn)
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	apiV1 := app.Group("/api/v1")

	api.UsersRoute(apiV1.Group("/users"), usersController)

	apiV1.Route("/urls", api.UrlsRoute)

	return app
}
