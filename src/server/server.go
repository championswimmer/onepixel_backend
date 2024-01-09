package server

import (
	"fmt"
	"onepixel_backend/src/config"
	"onepixel_backend/src/docs"
	_ "onepixel_backend/src/docs"
	"onepixel_backend/src/routes/api"
	"onepixel_backend/src/server/logger"

	"onepixel_backend/src/routes/redirect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

const (
	ENVIRONMENT_NAME_PRODUCTION = "production"
)

// CreateAdminApp creates the fiber app
//
//	@title		onepixel API
//	@version	0.1
//	@description.markdown
//	@termsOfService				https://github.com/championswimmer/onepixel_backend
//	@contact.name				Arnav Gupta
//	@contact.email				dev@championswimmer.in
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@host						onepixel.link
//	@BasePath					/api/v1
//	@schemes					http https
//
//	@tag.name					users
//	@tag.description			Operations about users
//	@tag.name					urls
//	@tag.description			Operations about urls
//
//	@securityDefinitions.apiKey	BearerToken
//	@in							header
//	@name						Authorization
//	@securityDefinitions.apiKey	APIKeyAuth
//	@in							header
//	@name						X-API-Key
func CreateAdminApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	switch config.Env {
	case ENVIRONMENT_NAME_PRODUCTION:
		docs.SwaggerInfo.Host = config.AdminHost
	default:
		app.Use(logger.NewLogger())
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.AdminHost, config.Port)
	}

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute(db))
	apiV1.Route("/urls", api.UrlsRoute(db))

	app.Get("/docs/*", swagger.HandlerDefault)

	return app
}

func CreateMainApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Route("/", redirect.RedirectRoute(db))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	return app
}
