package redirect

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/db/models"
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
		router.Get("/:group/:shortcode", redirectGroupedShortCode)
		router.Get("/:shortcode", redirectShortCode)
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
		return sendRedirectError(ctx, urlErr)
	}

	return redirectToDestination(ctx, url, nil)
}

func redirectGroupedShortCode(ctx *fiber.Ctx) error {
	group := ctx.Params("group")
	validErr := validators.ValidateRedirectGroupRequest(group)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	shortcode := ctx.Params("shortcode")
	validErr = validators.ValidateRedirectShortCodeRequest(shortcode)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	urlGroup, groupErr := urlsController.GetUrlGroupByShortPath(group)
	if groupErr != nil {
		return sendRedirectError(ctx, groupErr)
	}

	url, urlErr := urlsController.GetUrlWithShortCodeInGroup(shortcode, urlGroup.ID)
	if urlErr != nil {
		return sendRedirectError(ctx, urlErr)
	}

	return redirectToDestination(ctx, url, urlGroup)
}

func redirectToDestination(ctx *fiber.Ctx, url *models.Url, urlGroup *models.UrlGroup) error {
	canonicalShortURL := controllers.CanonicalShortURL(urlGroup, url.ShortURL)
	eventsController.LogRedirectAsync(&controllers.EventRedirectData{
		ShortUrlID: url.ID,
		UrlGroupID: url.UrlGroupID,
		ShortURL:   canonicalShortURL,
		CreatorID:  url.CreatorID,
		IPAddress:  strings.Split(ctx.Get("X-Forwarded-For"), ",")[0],
		UserAgent:  ctx.Get("User-Agent"),
		Referer:    ctx.Get("Referer"),
	})

	// cache for 1 min only
	// ctx.Response().Header.Set("Cache-Control", "public, max-age=60")
	return ctx.Redirect(url.LongURL, fiber.StatusMovedPermanently)
}

func sendRedirectError(ctx *fiber.Ctx, err error) error {
	var e *controllers.UrlError
	if errors.As(err, &e) {
		status, message := e.ErrorDetails()
		return ctx.Status(status).JSON(dtos.CreateErrorResponse(status, message))
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, err.Error()))
}
