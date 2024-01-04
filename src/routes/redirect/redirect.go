package redirect

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"onepixel_backend/src/controllers"
)

var urlsController *controllers.UrlsController

func RedirectRoute(db *gorm.DB) func(router fiber.Router) {
	urlsController = controllers.CreateUrlsController(db)
	return func(router fiber.Router) {
		router.Get("/:shortcode", redirectShortCode)
		router.Get("/:group/:shortcode", redirectGroupedShortCode)
	}

}

func redirectShortCode(ctx *fiber.Ctx) error {
	return ctx.SendString("Redirect")
}

func redirectGroupedShortCode(ctx *fiber.Ctx) error {
	return ctx.SendString("Redirect")
}
