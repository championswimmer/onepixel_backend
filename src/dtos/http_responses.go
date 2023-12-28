package dtos

import "onepixel_backend/src/models"

type UserResponse struct {
	ID    uint64  `json:"id" example:"1"`
	Email string  `json:"email" example:"user@test.com"`
	Token *string `json:"token" example:"<JWT_TOKEN>"`
}

type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Something went wrong"`
}

func CreateUserResponseFromUser(user *models.User, token *string) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}
}

func CreateErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Message: message,
	}
}
