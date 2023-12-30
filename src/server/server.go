package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
	_ "onepixel_backend/src/docs"
	"onepixel_backend/src/routes/api"
)

// CreateApp creates the fiber app
//
//	@title						onepixel API
//	@version					0.1
//	@description				1px.li URL Shortner API
//	@termsOfService				https://github.com/championswimmer/onepixel_backend
//	@contact.name				Arnav Gupta
//	@contact.email				dev@championswimmer.in
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@host						onepixel.link
//	@BasePath					/api/v1
//	@schemes					https
//	@securityDefinitions.apiKey	BearerToken
//	@in							header
//	@name						Authorization
//	@securityDefinitions.apiKey	APIKeyAuth
//	@in							header
//	@name						X-API-Key
//	@security					APIKeyAuth
func CreateApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute(db))
	apiV1.Route("/urls", api.UrlsRoute(db))

	app.Get("/docs/*", swagger.HandlerDefault)

	return app
}
