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
	// Splice out the "Bearer " prefix, if it exists
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		authHeader = authHeader[7:]
	}
	user, err := ValidateJWT(authHeader)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid JWT token" + err.Error(),
		})
	}
	c.Locals(config.LOCALS_USER, user)
	return c.Next()
}

func OptionalJwtAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}
	// Splice out the "Bearer " prefix, if it exists
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		authHeader = authHeader[7:]
	}
	user, err := ValidateJWT(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid JWT token",
		})
	}
	c.Locals(config.LOCALS_USER, user)
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
	c.Locals(config.LOCALS_ADMIN, true)
	return c.Next()
}

func OptionalAdminApiKeyAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("X-API-Key")
	if authHeader == "" {
		return c.Next()
	}
	if authHeader != config.AdminApiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid X-API-Key",
		})
	}
	c.Locals(config.LOCALS_ADMIN, true)
	return c.Next()
}
