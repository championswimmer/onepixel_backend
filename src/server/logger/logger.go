package logger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
		TimeFormat: "Mon, 02 Jan 2006 15:04:05 GMT",
	})
}
