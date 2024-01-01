package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"onepixel_backend/src/config"
	"onepixel_backend/src/db"
	"onepixel_backend/src/server"
	"onepixel_backend/src/utils/applogger"
	"strings"
)

func main() {
	// Initialize the database
	db := lo.Must(db.GetDB())

	// Create the app
	adminApp := server.CreateAdminApp(db)
	mainApp := server.CreateMainApp(db)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		host := strings.Split(c.Hostname(), ":")[0]
		switch host {
		case config.AdminHost:
			adminApp.Handler()(c.Context())
			return nil
		case config.MainHost:
			mainApp.Handler()(c.Context())
			return nil
		default:
			c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "called via unsupported host",
			})
			return nil
		}
	})

	applogger.Fatal(app.Listen(fmt.Sprintf(":%s", config.Port)))
}
