package security

import (
	"github.com/gofiber/fiber/v2"
	"onepixel_backend/src/config"
)

// MandatoryJwtAuthMiddleware makes authentication mandatory
// will return 401 if no Authorization header is provided or if the JWT is invalid
// saves the user in the context locals as "user"
func MandatoryJwtAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: No Authorization header provided",
		})
	}
	user, err := ValidateJWT(authHeader)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid JWT token" + err.Error(),
		})
	}
	c.Locals("user", user)
	return c.Next()
}

func OptionalJwtAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}
	user, err := ValidateJWT(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid JWT token",
		})
	}
	c.Locals("user", user)
	return c.Next()
}

func MandatoryAdminApiKeyAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("X-API-Key")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: No X-API-Key header provided",
		})
	}
	if authHeader != config.AdminApiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid X-API-Key",
		})
	}
	c.Locals("admin", true)
	return c.Next()
}
