package dtos

import "onepixel_backend/src/models"

type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Email string `json:"email" example:"user@test.com"`
}

type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Something went wrong"`
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
