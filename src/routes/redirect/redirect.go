package redirect

import (
	"errors"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/server/validators"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

//var urlsController *controllers.UrlsController

var urlsController models.IUrlController

func RedirectRoute(db *gorm.DB) func(router fiber.Router) {
	urlsController = controllers.CreateUrlsController(db)
	return func(router fiber.Router) {
		router.Get("/:shortcode", redirectShortCode)
		router.Get("/:group/:shortcode", redirectGroupedShortCode)
	}

}

func redirectShortCode(ctx *fiber.Ctx) error {
	shortcode := ctx.Params("shortcode")
	validErr := validators.ValidateRedirectShortCodeRequest(shortcode)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	url, urlErr := urlsController.GetUrlWithShortCode(shortcode)
	if urlErr != nil {
		var e *controllers.UrlError
		if errors.As(urlErr, &e) {
			return ctx.Status(fiber.StatusNotFound).JSON(dtos.CreateErrorResponse(e.ErrorDetails()))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, urlErr.Error()))
	}
	return ctx.Redirect(url.LongURL, fiber.StatusMovedPermanently)
}

func redirectGroupedShortCode(ctx *fiber.Ctx) error {
	return ctx.SendString("Redirect")
}
