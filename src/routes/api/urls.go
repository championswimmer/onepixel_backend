package api

import "github.com/gofiber/fiber/v2"

// UrlsRoute /api/v1/urls
func UrlsRoute(router fiber.Router) {
	router.Get("/", getAllUrls)
}

func getAllUrls(ctx *fiber.Ctx) error {
	return ctx.SendString("GetAllUsers")
}
