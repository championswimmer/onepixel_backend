package validators

import (
	"github.com/gofiber/fiber/v2"
	"onepixel_backend/src/dtos"
)

var mandatoryUrlDtoFieldError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "long_url is required",
}

func ValidateCreateUrlRequest(dto *dtos.CreateUrlRequest) *ValidationError {
	if dto.LongUrl == "" {
		return mandatoryUrlDtoFieldError
	}
	return nil
}
