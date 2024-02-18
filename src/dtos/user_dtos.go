package dtos

import "onepixel_backend/src/db/models"

/// Requests

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/// Responses

type UserResponse struct {
	ID    uint64  `json:"id" example:"1"`
	Email string  `json:"email" example:"user@test.com"`
	Token *string `json:"token" example:"<JWT_TOKEN>"`
}

/// Converters

func CreateUserResponseFromUser(user *models.User, token *string) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}
}
