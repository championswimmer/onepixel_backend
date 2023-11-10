package dtos

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUrlRequest struct {
	LongUrl string `json:"longUrl"`
	GroupId uint   `json:"groupId"`
}
