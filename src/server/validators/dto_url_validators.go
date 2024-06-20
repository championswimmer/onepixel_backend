package validators

import (
	"github.com/gofiber/fiber/v2"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/utils"
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
		message: "Shortcode exceeds the maximum allowed length of 10 characters",
	}

	ShortcodeEmptyError = &ValidationError{
		status:  fiber.ErrBadRequest.Code,
		message: "Shortcode is empty",
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

func ValidateSpecificShortCodeRequest(shortcode string) *ValidationError {
	if shortcode == "" {
		return ShortcodeEmptyError
	}
	if len(shortcode) > utils.MaxSafeStringLength {
		return ShortcodeTooLongError
	}
	return nil
}
