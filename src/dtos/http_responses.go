package dtos

import "onepixel_backend/src/models"

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func UserResponseFromUser(user *models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}
