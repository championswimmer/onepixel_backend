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
		router.Get("/", security.OptionalJwtAuthMiddleware, security.OptionalAdminApiKeyAuthMiddleware, getAllUrls)
		router.Get("/groups", security.MandatoryJwtAuthMiddleware, getUrlGroups)
		router.Post("/groups", security.MandatoryAdminApiKeyAuthMiddleware, createUrlGroup)
		router.Post("/groups/:group/shorten", security.MandatoryJwtAuthMiddleware, createGroupedRandomUrl)
		router.Post("/groups/:group/shorten/:shortcode", security.MandatoryJwtAuthMiddleware, createGroupedSpecificUrl)
		router.Get("/groups/:group/:shortcode", getGroupedUrlInfo)
		router.Post("/", security.MandatoryJwtAuthMiddleware, createRandomUrl)
		router.Put("/:shortcode", security.MandatoryJwtAuthMiddleware, createSpecificUrl)
		router.Get("/:shortcode", getUrlInfo)
	}
}

// getAllUrls
//
//	@Summary		Get all urls
//	@Description	Get all urls
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dtos.UrlResponse
//	@Failure		500	{object}	dtos.ErrorResponse	"something went wrong"
//	@Router			/urls [get]
//	@Security		BearerToken
func getAllUrls(ctx *fiber.Ctx) error {
	userObj := ctx.Locals(config.LOCALS_USER)
	admin := ctx.Locals(config.LOCALS_ADMIN)
	if userObj == nil && admin == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dtos.CreateErrorResponse(fiber.StatusUnauthorized, "unauthorised: neither admin key nor auth header provided"))
	}

	var urls []models.Url
	if admin != nil && admin.(bool) {
		u := ctx.Query("userid")
		var (
			userID *uint64
			err    error
		)
		if u != "" {
			uID, err := strconv.ParseUint(u, 10, 64)
			if err != nil {
				return ctx.Status(fiber.StatusNotFound).JSON(dtos.CreateErrorResponse(fiber.StatusNotFound, "invalid user id"))
			}
			userID = &uID
		}

		urls, err = urlsController.GetUrls(userID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "something went wrong"))
		}
	}

	// for user who is not admin
	if userObj != nil && admin == nil {
		var err error
		user := *userObj.(*models.User)
		urls, err = urlsController.GetUrls(&user.ID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "something went wrong"))
		}
	}

	urlResponses := make([]dtos.UrlResponse, 0)
	for _, url := range urls {
		urlResponses = append(urlResponses, dtos.CreateUrlResponse(&url))
	}

	return ctx.Status(fiber.StatusOK).JSON(urlResponses)
}

// getUrlGroups
//
//	@Summary		Get URL groups for the current user
//	@Description	Get all URL groups created by the authenticated user
//	@ID				get-url-groups
//	@Tags			urls
//	@Produce		json
//	@Success		200	{array}		dtos.UrlGroupResponse
//	@Failure		500	{object}	dtos.ErrorResponse	"something went wrong"
//	@Router			/urls/groups [get]
//	@Security		BearerToken
func getUrlGroups(ctx *fiber.Ctx) error {
	user := ctx.Locals(config.LOCALS_USER).(*models.User)

	groups, err := urlsController.GetUrlGroupsByCreator(user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "something went wrong"))
	}

	groupResponses := make([]dtos.UrlGroupResponse, 0)
	for _, group := range groups {
		groupResponses = append(groupResponses, dtos.CreateUrlGroupResponse(&group))
	}

	return ctx.Status(fiber.StatusOK).JSON(groupResponses)
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
	validErr = validators.ValidateCreateShortCodeRequest(shortcode)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	createdUrl, createErr := urlsController.CreateSpecificShortUrl(shortcode, cur.LongUrl, user.ID)
	if createErr != nil {
		return sendUrlControllerError(ctx, createErr)
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUrlResponse(createdUrl))
}

// createUrlGroup
//
//	@Summary		Create URL group
//	@Description	Create a URL group for a specific user
//	@ID				create-url-group
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			urlGroup	body		dtos.CreateUrlGroupRequest	true	"URL group"
//	@Success		201			{object}	dtos.UrlGroupResponse
//	@Failure		400			{object}	dtos.ErrorResponse	"The request body is not valid"
//	@Failure		409			{object}	dtos.ErrorResponse	"URL group already exists"
//	@Failure		422			{object}	dtos.ErrorResponse	"invalid group"
//	@Router			/urls/groups [post]
//	@Security		APIKeyAuth
func createUrlGroup(ctx *fiber.Ctx) error {
	cur, parseErr := parsers.ParseBody[dtos.CreateUrlGroupRequest](ctx)
	if parseErr != nil {
		return parsers.SendParsingError(ctx, parseErr)
	}

	validErr := validators.ValidateCreateUrlGroupDtoRequest(cur)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	urlGroup, createErr := urlsController.CreateUrlGroup(cur.ShortPath, cur.CreatorID)
	if createErr != nil {
		return sendUrlControllerError(ctx, createErr)
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUrlGroupResponse(urlGroup))
}

// createGroupedSpecificUrl
//
//	@Summary		Create grouped short URL with specific shortcode
//	@Description	Create a grouped short URL for a URL group owned by the authenticated user
//	@ID				create-grouped-specific-url
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			group		path		string					true	"URL group"
//	@Param			shortcode	path		string					true	"Shortcode"
//	@Param			url			body		dtos.CreateUrlRequest	true	"Url"
//	@Success		201			{object}	dtos.UrlResponse
//	@Failure		400			{object}	dtos.ErrorResponse	"The request body is not valid"
//	@Failure		403			{object}	dtos.ErrorResponse	"URL group does not belong to the user"
//	@Failure		404			{object}	dtos.ErrorResponse	"URL group not found"
//	@Failure		409			{object}	dtos.ErrorResponse	"URL already exists"
//	@Failure		422			{object}	dtos.ErrorResponse	"invalid shortcode"
//	@Router			/urls/groups/{group}/shorten/{shortcode} [post]
//	@Security		BearerToken
func createGroupedSpecificUrl(ctx *fiber.Ctx) error {
	user := ctx.Locals(config.LOCALS_USER).(*models.User)
	urlGroup, handled, groupErr := getOwnedUrlGroup(ctx, user)
	if groupErr != nil {
		return groupErr
	}
	if handled {
		return nil
	}

	cur, parseErr := parsers.ParseBody[dtos.CreateUrlRequest](ctx)
	if parseErr != nil {
		return parsers.SendParsingError(ctx, parseErr)
	}

	shortcode := ctx.Params("shortcode")
	validErr := validators.ValidateCreateUrlRequest(cur)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	validErr = validators.ValidateCreateShortCodeRequest(shortcode)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	createdUrl, createErr := urlsController.CreateSpecificShortUrlInGroup(shortcode, cur.LongUrl, user.ID, urlGroup.ID)
	if createErr != nil {
		return sendUrlControllerError(ctx, createErr)
	}

	createdUrl.UrlGroup = *urlGroup
	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUrlResponse(createdUrl))
}

// createGroupedRandomUrl
//
//	@Summary		Create grouped short URL with random shortcode
//	@Description	Create a grouped short URL for a URL group owned by the authenticated user
//	@ID				create-grouped-random-url
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			group	path		string					true	"URL group"
//	@Param			url		body		dtos.CreateUrlRequest	true	"Url"
//	@Success		201		{object}	dtos.UrlResponse
//	@Failure		400		{object}	dtos.ErrorResponse	"The request body is not valid"
//	@Failure		403		{object}	dtos.ErrorResponse	"URL group does not belong to the user"
//	@Failure		404		{object}	dtos.ErrorResponse	"URL group not found"
//	@Failure		422		{object}	dtos.ErrorResponse	"invalid group"
//	@Router			/urls/groups/{group}/shorten [post]
//	@Security		BearerToken
func createGroupedRandomUrl(ctx *fiber.Ctx) error {
	user := ctx.Locals(config.LOCALS_USER).(*models.User)
	urlGroup, handled, groupErr := getOwnedUrlGroup(ctx, user)
	if groupErr != nil {
		return groupErr
	}
	if handled {
		return nil
	}

	cur, parseErr := parsers.ParseBody[dtos.CreateUrlRequest](ctx)
	if parseErr != nil {
		return parsers.SendParsingError(ctx, parseErr)
	}

	validErr := validators.ValidateCreateUrlRequest(cur)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	createdUrl, createErr := urlsController.CreateRandomShortUrlInGroup(cur.LongUrl, user.ID, urlGroup.ID)
	if createErr != nil {
		return sendUrlControllerError(ctx, createErr)
	}

	createdUrl.UrlGroup = *urlGroup
	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUrlResponse(createdUrl))
}

// getGroupedUrlInfo
//
//	@Summary		Get grouped URL info
//	@Description	Get URL info for a URL in a specific group
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			group		path		string	true	"URL group"
//	@Param			shortcode	path		string	true	"Shortcode"
//	@Success		200			{object}	dtos.UrlInfoResponse
//	@Failure		404			{object}	dtos.ErrorResponse	"URL or group not found"
//	@Router			/urls/groups/{group}/{shortcode} [get]
func getGroupedUrlInfo(ctx *fiber.Ctx) error {
	group := ctx.Params("group")
	shortcode := ctx.Params("shortcode")

	urlGroup, groupErr := urlsController.GetUrlGroupByShortPath(group)
	if groupErr != nil {
		return sendUrlControllerError(ctx, groupErr)
	}

	longUrl, hitCount, err := urlsController.GetUrlInfoInGroup(shortcode, urlGroup.ID)
	if err != nil {
		return sendUrlControllerError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.CreateUrlInfoResponse(longUrl, hitCount))
}

// getUrlInfo
//
//	@Summary		Get URL info
//	@Description	Get URL info
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			shortcode	path		string	true	"Shortcode"
//	@Success		200			{object}	dtos.UrlInfoResponse
//	@Failure		404			{object}	dtos.ErrorResponse	"URL not found"
//	@Router			/urls/{shortcode} [get]
func getUrlInfo(ctx *fiber.Ctx) error {
	shortcode := ctx.Params("shortcode")
	validErr := validators.ValidateCreateShortCodeRequest(shortcode)
	if validErr != nil {
		return validators.SendValidationError(ctx, validErr)
	}

	longUrl, hitCount, err := urlsController.GetUrlInfo(shortcode)
	if err != nil {
		return sendUrlControllerError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.CreateUrlInfoResponse(longUrl, hitCount))
}

func getOwnedUrlGroup(ctx *fiber.Ctx, user *models.User) (*models.UrlGroup, bool, error) {
	group := ctx.Params("group")
	validErr := validators.ValidateCreateUrlGroupRequest(group)
	if validErr != nil {
		return nil, true, validators.SendValidationError(ctx, validErr)
	}

	urlGroup, groupErr := urlsController.GetUrlGroupByShortPath(group)
	if groupErr != nil {
		return nil, true, sendUrlControllerError(ctx, groupErr)
	}
	if urlGroup.CreatorID != user.ID {
		return nil, true, ctx.Status(fiber.StatusForbidden).JSON(dtos.CreateErrorResponse(controllers.UrlGroupForbiddenError.ErrorDetails()))
	}

	return urlGroup, false, nil
}

func sendUrlControllerError(ctx *fiber.Ctx, err error) error {
	var e *controllers.UrlError
	if errors.As(err, &e) {
		status, message := e.ErrorDetails()
		return ctx.Status(status).JSON(dtos.CreateErrorResponse(status, message))
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "something went wrong"))
}
