package server

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/models"
	"onepixel_backend/src/routes/api"
)

func CreateApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("/testJWT", func(ctx *fiber.Ctx) error {
		var user models.User
		db.Last(&user)
		token := auth.CreateJWTFromUser(&user)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	})

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute(db))
	apiV1.Route("/urls", api.UrlsRoute(db))

	return app
}
