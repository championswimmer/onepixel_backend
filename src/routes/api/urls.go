package api

import (
	"errors"
	"onepixel_backend/src/config"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/security"
	"onepixel_backend/src/server/parsers"
	"onepixel_backend/src/server/validators"

	"github.com/gofiber/fiber/v2"
)

var urlsController *controllers.UrlsController

// UrlsRoute
func UrlsRoute() func(router fiber.Router) {
	// initialize UrlsController
	urlsController = controllers.CreateUrlsController()

	return func(router fiber.Router) {
		router.Get("/", getAllUrls)
		router.Post("/", security.MandatoryJwtAuthMiddleware, createRandomUrl)
		router.Put("/:shortcode", security.MandatoryJwtAuthMiddleware, createSpecificUrl)
	}
}

// getAllUrls
//
//	@Summary		Get all urls
//	@Description	Get all urls
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"GetAllUsers"
//	@Router			/urls [get]
//	@security		BearerToken
func getAllUrls(ctx *fiber.Ctx) error {
	return ctx.SendString("GetAllUsers")
}

// createRandomUrl
//
//	@Summary		Create random short url
//	@Description	Create random short url
//	@ID				create-random-url
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			url	body		dtos.CreateUrlRequest	true	"Url"
//	@Success		201	{object}	dtos.UrlResponse
//	@Failure		400	{object}	dtos.ErrorResponse	"The request body is not valid"
//	@Failure		422	{object}	dtos.ErrorResponse	"long_url is required to create url"
//	@Router			/urls [post]
//	@Security		BearerToken
func createRandomUrl(ctx *fiber.Ctx) error {
	user := ctx.Locals(config.LOCALS_USER).(*models.User)

	cur, parseErr := parsers.ParseBody[dtos.CreateUrlRequest](ctx)
	if parseErr != nil {
		return parsers.SendParsingError(ctx, parseErr)
	}

	validErr := validators.ValidateCreateUrlRequest(cur)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	createdUrl, createErr := urlsController.CreateRandomShortUrl(cur.LongUrl, user.ID)
	if createErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "something went wrong"))
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUrlResponse(createdUrl))
}

// createSpecificUrl
//
//	@Summary		Create specific short url
//	@Description	Create specific short url
//	@ID				create-specific-url
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			shortcode	path		string					true	"Shortcode"
//	@Param			url			body		dtos.CreateUrlRequest	true	"Url"
//	@Success		201			{object}	dtos.UrlResponse
//	@Failure		400			{object}	dtos.ErrorResponse	"The request body is not valid"
//	@Failure		422			{object}	dtos.ErrorResponse	"long_url is required to create url"
//	@Failure		409			{object}	dtos.ErrorResponse	"Shortcode already exists"
//	@Failure		403			{object}	dtos.ErrorResponse	"Shortcode is not allowed"
//	@Router			/urls/{shortcode} [put]
//	@Security		BearerToken
func createSpecificUrl(ctx *fiber.Ctx) error {
	user := ctx.Locals(config.LOCALS_USER).(*models.User)

	cur, parseErr := parsers.ParseBody[dtos.CreateUrlRequest](ctx)
	if parseErr != nil {
		return parsers.SendParsingError(ctx, parseErr)
	}

	validErr := validators.ValidateCreateUrlRequest(cur)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	shortcode := ctx.Params("shortcode")
	if shortcode == "" {
		// TODO: handle unacceptable/reserved shortcodes properly in controller
		panic("shortcode is empty")
	}

	createdUrl, createErr := urlsController.CreateSpecificShortUrl(shortcode, cur.LongUrl, user.ID)
	if createErr != nil {
		var e *controllers.UrlError
		if errors.As(createErr, &e) {
			if errors.Is(e, controllers.UrlExistsError) {
				return ctx.Status(fiber.StatusConflict).JSON(dtos.CreateErrorResponse(e.ErrorDetails()))
			}
			if errors.Is(e, controllers.UrlForbiddenError) {
				return ctx.Status(fiber.StatusForbidden).JSON(dtos.CreateErrorResponse(e.ErrorDetails()))
			}
		} else {
			return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "something went wrong"))
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUrlResponse(createdUrl))
}

func createGroupedRandomUrl(ctx *fiber.Ctx) error {
	return ctx.SendString("createGroupedRandomUrl")
}

func createGroupedSpecificUrl(ctx *fiber.Ctx) error {
	return ctx.SendString("createGroupedSpecificUrl")
}
