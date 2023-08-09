package models

// Create User request body

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User response body

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}
