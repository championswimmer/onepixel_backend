package api

import (
	"onepixel_backend/src/config"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/security"
	"onepixel_backend/src/utils/applogger"

	"github.com/gofiber/fiber/v2"
)

var eventsController *controllers.EventsController

// StatsRoute
func StatsRoute() func(router fiber.Router) {
	// initialize EventsController
	eventsController = controllers.CreateEventsController()

	return func(router fiber.Router) {
		router.Get("/", security.OptionalJwtAuthMiddleware, getAllStats)
		router.Get("/:shortcode" /*security.MandatoryJwtAuthMiddleware,*/, getStatsForShortCode)
		// TODO: add stats for grouped shortcodes
	}
}

func getAllStats(ctx *fiber.Ctx) error {
	// TODO: handle null case
	user := ctx.Locals(config.LOCALS_USER).(*models.User)
	stats, err := eventsController.GetRedirectsCountForUserId(user.ID)

	if err != nil {
		applogger.Error(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.CreateErrorResponse(fiber.StatusInternalServerError, "something went wrong"))
	}

	return ctx.Status(fiber.StatusOK).JSON(stats)
}

func getStatsForShortCode(ctx *fiber.Ctx) error {
	// stats := eventsController.GetRedirectsCountForUserId("")
	// applogger.Info("Stats: ", len(stats))
	return ctx.SendString("GetStats")
}
