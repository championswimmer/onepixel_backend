package validators

import (
	"github.com/gofiber/fiber/v2"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/utils"
)

var mandatoryUrlDtoFieldError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "long_url is required",
}

var invalidShortCodeError = &ValidationError{
	status:  fiber.StatusNotFound,
	message: "Invalid short code",
}

func ValidateCreateUrlRequest(dto *dtos.CreateUrlRequest) *ValidationError {
	if dto.LongUrl == "" {
		return mandatoryUrlDtoFieldError
	}
	return nil
}

func ValidateRedirectShortCodeRequest(shortcode string) *ValidationError {
	if shortcode == "" {
		return invalidShortCodeError
	}
	if len(shortcode) > utils.MaxSafeStringLength {
		return invalidShortCodeError
	}
	return nil
}
