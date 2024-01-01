package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
	"onepixel_backend/src/config"
	"onepixel_backend/src/docs"
	_ "onepixel_backend/src/docs"
	"onepixel_backend/src/routes/api"
)

// CreateAdminApp creates the fiber app
//
//	@title						onepixel API
//	@version					0.1
//	@description				1px.li URL Shortner API
//	@termsOfService				https://github.com/championswimmer/onepixel_backend
//	@contact.name				Arnav Gupta
//	@contact.email				dev@championswimmer.in
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@host						api.onepixel.link
//	@BasePath					/api/v1
//	@schemes					http https
//	@securityDefinitions.apiKey	BearerToken
//	@in							header
//	@name						Authorization
//	@securityDefinitions.apiKey	APIKeyAuth
//	@in							header
//	@name						X-API-Key
//	@security					APIKeyAuth
func CreateAdminApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute(db))
	apiV1.Route("/urls", api.UrlsRoute(db))

	if config.Env == "production" {
		docs.SwaggerInfo.Host = config.AdminHost
	} else {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.AdminHost, config.Port)
	}

	app.Get("/docs/*", swagger.HandlerDefault)

	return app
}

func CreateMainApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	return app
}
