package validators

import (
	"github.com/gofiber/fiber/v2"
	"onepixel_backend/src/dtos"
)

var mandatoryUserDtoFieldsError = &ValidationError{
	status:  fiber.StatusUnprocessableEntity,
	message: "Email and password are required",
}

func ValidateCreateUserRequest(dto *dtos.CreateUserRequest) *ValidationError {
	if dto.Email == "" || dto.Password == "" {
		return mandatoryUserDtoFieldsError
	}
	return nil
}

func ValidateLoginUserRequest(dto *dtos.LoginUserRequest) *ValidationError {
	if dto.Email == "" || dto.Password == "" {
		return mandatoryUserDtoFieldsError
	}
	return nil
}
