package dtos

import "onepixel_backend/src/models"

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

func UserResponseFromUser(user *models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

func GetErrorResponse(status uint, message string) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Message: message,
	}
}
