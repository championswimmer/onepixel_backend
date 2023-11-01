package dtos

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
