package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"onepixel_backend/src/config"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/docs"
	_ "onepixel_backend/src/docs"
	"onepixel_backend/src/routes/api"
	"onepixel_backend/src/routes/redirect"
	"onepixel_backend/src/utils/applogger"
	"strings"
	"time"
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
func CreateAdminApp() *fiber.App {
	app := fiber.New()

	apiV1 := app.Group("/api/v1")

	apiV1.Route("/users", api.UsersRoute())
	apiV1.Route("/urls", api.UrlsRoute())
	apiV1.Route("/stats", api.StatsRoute())

	if config.Env == "production" {
		docs.SwaggerInfo.Host = config.AdminHost
	} else {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.AdminHost, config.Port)
	}

	app.Get("/docs/*", swagger.HandlerDefault)
	app.Static("/", "./public_html", fiber.Static{
		Compress:      true,
		MaxAge:        60 * 60, // 1 hour
		CacheDuration: time.Hour * 1,
	})
	return app
}

func CreateMainApp() *fiber.App {
	app := fiber.New()
	eventsController := controllers.CreateEventsController()

	app.Route("/", redirect.RedirectRoute())
	app.Get("/", func(c *fiber.Ctx) error {
		redirPath := c.Protocol() + "://" + config.AdminHost
		if config.Env == "local" {
			redirPath += ":" + config.Port
		}
		applogger.Info("redirect: root: " + c.OriginalURL())
		eventsController.LogRedirectAsync(&controllers.EventRedirectData{
			UrlGroupID: 0,
			ShortURL:   "/",
			CreatorID:  0,
			IPAddress:  strings.Split(c.Get("X-Forwarded-For"), ",")[0],
			UserAgent:  c.Get("User-Agent"),
			Referer:    c.Get("Referer"),
		})
		return c.Redirect(redirPath, fiber.StatusMovedPermanently)
	})

	return app
}
