package dtos

/// User requests

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/// URL Requests

type CreateUrlRequest struct {
	LongUrl string `json:"long_url"`
}
