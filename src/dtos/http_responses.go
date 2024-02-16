package dtos

import "onepixel_backend/src/db/models"

type AppError interface {
	ErrorDetails() (int, string)
}

type UserResponse struct {
	ID    uint64  `json:"id" example:"1"`
	Email string  `json:"email" example:"user@test.com"`
	Token *string `json:"token" example:"<JWT_TOKEN>"`
}

type UrlResponse struct {
	ShortURL  string `json:"short_url" example:"nhg145"`
	LongURL   string `json:"long_url" example:"https://www.google.com"`
	CreatorID uint64 `json:"creator_id" example:"1"`
}

type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Something went wrong"`
}

type RedirectCountStatsResponse struct {
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

func CreateUrlResponse(url *models.Url) UrlResponse {
	return UrlResponse{
		ShortURL:  url.ShortURL,
		LongURL:   url.LongURL,
		CreatorID: url.CreatorID,
	}
}
