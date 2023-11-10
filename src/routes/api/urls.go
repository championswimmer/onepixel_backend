package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"onepixel_backend/src/auth"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/models"
)

var urlController *controllers.UrlController

// UrlsRoute /api/v1/urls
func UrlsRoute(db *gorm.DB) func(router fiber.Router) {
	urlController = controllers.NewUrlController(db)
	return func(router fiber.Router) {

		router.Get("/", getAllUrls)
		router.Post("/", auth.MandatoryAuthMiddleware, createShortUrl)
	}
}

func getAllUrls(ctx *fiber.Ctx) error {
	return ctx.SendString("GetAllUsers")
}
func createShortUrl(ctx *fiber.Ctx) error {
	var url = new(dtos.CreateUrlRequest)
	lo.Must0(ctx.BodyParser(url))
	userIdInterface := ctx.Locals("user")
	userID, ok := userIdInterface.(*models.User)
	if !ok {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "User does not exist",
		})
	}
	savedUrl, err := urlController.Create(url.LongUrl, url.GroupId, userID.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Try Again",
			})
		} else {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}
	return ctx.Status(fiber.StatusCreated).JSON(dtos.UrlResponseFromUrl(savedUrl))
}
