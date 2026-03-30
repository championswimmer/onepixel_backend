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

var mandatoryUrlGroupCreatorFieldError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "creator_id is required",
}

var invalidShortCodeError = &ValidationError{
	status:  fiber.StatusNotFound,
	message: "Invalid short code",
}

var invalidUrlGroupError = &ValidationError{
	status:  fiber.StatusNotFound,
	message: "Invalid url group",
}

var invalidCreateShortCodeError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "invalid shortcode",
}

var invalidCreateUrlGroupError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "invalid group",
}

func ValidateCreateUrlRequest(dto *dtos.CreateUrlRequest) *ValidationError {
	if dto.LongUrl == "" {
		return mandatoryUrlDtoFieldError
	}
	return nil
}

func ValidateCreateUrlGroupDtoRequest(dto *dtos.CreateUrlGroupRequest) *ValidationError {
	if dto.CreatorID == 0 {
		return mandatoryUrlGroupCreatorFieldError
	}
	return ValidateCreateUrlGroupRequest(dto.ShortPath)
}

func ValidateRedirectShortCodeRequest(shortcode string) *ValidationError {
	return validateRadix64Token(shortcode, invalidShortCodeError)
}

func ValidateRedirectGroupRequest(group string) *ValidationError {
	return validateRadix64Token(group, invalidUrlGroupError)
}

func ValidateCreateShortCodeRequest(shortcode string) *ValidationError {
	return validateRadix64Token(shortcode, invalidCreateShortCodeError)
}

func ValidateCreateUrlGroupRequest(group string) *ValidationError {
	return validateRadix64Token(group, invalidCreateUrlGroupError)
}

func validateRadix64Token(token string, err *ValidationError) *ValidationError {
	if token == "" {
		return err
	}
	if len(token) > utils.MaxSafeStringLength {
		return err
	}
	for _, char := range token {
		if _, ok := utils.AlphabetIndex[char]; !ok {
			return err
		}
	}
	return nil
}
