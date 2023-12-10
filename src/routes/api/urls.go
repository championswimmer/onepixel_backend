package api

import "github.com/gofiber/fiber/v2"

// UrlsRoute
func UrlsRoute(router fiber.Router) {
	router.Get("/", getAllUrls)
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
