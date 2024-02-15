package api

import (
	"github.com/gofiber/fiber/v2"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/utils/applogger"
)

var eventsController *controllers.EventsController

// StatsRoute
func StatsRoute() func(router fiber.Router) {
	// initialize EventsController
	eventsController = controllers.CreateEventsController()

	return func(router fiber.Router) {
		router.Get("/", getStats)
	}
}

func getStats(ctx *fiber.Ctx) error {
	stats := eventsController.GetRedirectsCountForUserId("")
	applogger.Info("Stats: ", len(stats))
	return ctx.SendString("GetStats")
}
