package dtos

import "onepixel_backend/src/models"

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func CreateUserResponseFromUser(user *models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

func CreateErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Message: message,
	}
}
