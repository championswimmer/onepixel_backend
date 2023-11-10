package dtos

import "onepixel_backend/src/models"

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func UserResponseFromUser(user *models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

func LoginResponseFromUser(token string) LoginResponse {
	return LoginResponse{
		Token: token,
	}
}

func ErrorResponseFromServer(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
	}
}
