package redirect

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/server/validators"
	"onepixel_backend/src/utils/applogger"
	"strings"
)

var urlsController *controllers.UrlsController
var eventsController *controllers.EventsController

func RedirectRoute() func(router fiber.Router) {
	urlsController = controllers.CreateUrlsController()
	eventsController = controllers.CreateEventsController()

	return func(router fiber.Router) {
		router.Get("/:shortcode", redirectShortCode)
		router.Get("/:group/:shortcode", redirectGroupedShortCode)
	}

}

func redirectShortCode(ctx *fiber.Ctx) error {
	shortcode := ctx.Params("shortcode")
	applogger.Info("redirect: shortcode: " + ctx.OriginalURL())
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
	eventsController.LogRedirectAsync(&controllers.EventRedirectData{
		ShortUrlID: url.ID,
		UrlGroupID: url.UrlGroupID,
		ShortURL:   url.ShortURL,
		CreatorID:  url.CreatorID,
		IPAddress:  strings.Split(ctx.Get("X-Forwarded-For"), ",")[0],
		UserAgent:  ctx.Get("User-Agent"),
		Referer:    ctx.Get("Referer"),
	})
	return ctx.Redirect(url.LongURL, fiber.StatusMovedPermanently)
}

func redirectGroupedShortCode(ctx *fiber.Ctx) error {
	return ctx.SendString("Redirect")
}
