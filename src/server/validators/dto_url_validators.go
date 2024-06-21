package validators

import (
	"fmt"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/utils"

	"github.com/gofiber/fiber/v2"
)

var (
	mandatoryUrlDtoFieldError = &ValidationError{
		status:  fiber.StatusUnprocessableEntity,
		message: "long_url is required",
	}

	invalidShortCodeError = &ValidationError{
		status:  fiber.StatusNotFound,
		message: "Invalid short code",
	}
)

var (
	ShortcodeTooLongError = &ValidationError{
		status:  fiber.ErrBadRequest.Code,
		message: fmt.Sprintf("Shortcode exceeds the maximum allowed length of %d characters", utils.MaxSafeStringLength),
	}

	ShortcodeEmptyError = &ValidationError{
		status:  fiber.ErrBadRequest.Code,
		message: "Shortcode cannot be empty",
	}
)

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

// Validates the request by user to create a specific short URL
func ValidateSpecificShortCodeRequest(shortcode string) *ValidationError {
	if shortcode == "" {
		return ShortcodeEmptyError
	}
	if len(shortcode) > utils.MaxSafeStringLength {
		return ShortcodeTooLongError
	}
	return nil
}
