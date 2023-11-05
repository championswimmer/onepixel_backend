package dtos

import "onepixel_backend/src/models"

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type UrlResponse struct {
	ID       uint         `json:"id"`
	LongUrl  string       `json:"longUrl"`
	ShortUrl string       `json:"shortUrl"`
	GroupId  uint         `json:"groupId"`
	Creator  UserResponse `json:"creator"`
}

func UserResponseFromUser(user *models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}

func UrlResponseFromUrl(url *models.Url) UrlResponse {
	return UrlResponse{
		ID:       url.ID,
		LongUrl:  url.LongURL,
		ShortUrl: url.ShortURL,
		GroupId:  url.GroupID,
		Creator:  UserResponseFromUser(&url.Creator),
	}
}
