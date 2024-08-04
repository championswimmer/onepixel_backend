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
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var urlsController *controllers.UrlsController

// UrlsRoute
func UrlsRoute() func(router fiber.Router) {
	// initialize UrlsController
	urlsController = controllers.CreateUrlsController()

	return func(router fiber.Router) {
		router.Get("/", security.OptionalJwtAuthMiddleware, getAllUrls)
		router.Post("/", security.MandatoryJwtAuthMiddleware, createRandomUrl)
		router.Put("/:shortcode", security.MandatoryJwtAuthMiddleware, createSpecificUrl)
	}
}

// getAllUrls
//
//	@Summary		Get all urls (for users or admins)
//	@Description	Get all urls for the current user or all URLs for admins
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			X-API-Key	header		string	false	"Admin API Key (optional)"
//	@Param			userid		query		string	false	"Filter URLs by user ID (admin only)"
//	@Success		200			{array}		dtos.UrlResponse
//	@Failure		400			{object}	dtos.ErrorResponse	"Invalid user ID"
//	@Failure		401			{object}	dtos.ErrorResponse	"Unauthorized"
//	@Failure		500			{object}	dtos.ErrorResponse	"Failed to fetch URLs"
//	@Router			/urls [get]
//	@Security		BearerToken
//	@Security		ApiKeyAuth
func getAllUrls(ctx *fiber.Ctx) error {
	apiKey := ctx.Get("X-API-Key")
	isAdmin := apiKey == config.AdminApiKey

	var userId *uint64
	var err error

	if isAdmin {
		userIdStr := ctx.Query("userid")
		if userIdStr != "" {
			parsedId, parseErr := strconv.ParseUint(userIdStr, 10, 64)
			if parseErr != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(dtos.CreateErrorResponse(fiber.StatusBadRequest, "Invalid user ID"))
			}
			userId = &parsedId
		}
	} else {
		user, ok := ctx.Locals(config.LOCALS_USER).(*models.User)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(dtos.CreateErrorResponse(fiber.StatusUnauthorized, "Unauthorized"))
		}
		userId = &user.ID
	}

	var urls []models.Url
	if isAdmin {
		urls, err = urlsController.GetAllUrls(userId)
	} else {
		urls, err = urlsController.GetUrlsByUserId(*userId)
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "Failed to fetch URLs"))
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.CreateUrlsResponse(urls))
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
