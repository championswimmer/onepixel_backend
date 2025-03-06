package dtos

import (
	"onepixel_backend/src/config"
	"onepixel_backend/src/db/models"
)

/// Requests

type CreateUrlRequest struct {
	LongUrl string `json:"long_url"`
}

/// Responses

type UrlResponse struct {
	ShortURL  string `json:"short_url" example:"https://1px.li/nhg145"` // Example url will pick up host and protocol(http/https) based on the env
	LongURL   string `json:"long_url" example:"https://www.google.com"`
	CreatorID uint64 `json:"creator_id" example:"1"`
}

type UrlInfoResponse struct {
	LongURL  string `json:"long_url"`
	HitCount int64  `json:"hit_count"`
}

/// Converters

func CreateUrlResponse(url *models.Url) UrlResponse {
	return UrlResponse{
		ShortURL:  config.RedirUrlBase + url.ShortURL,
		LongURL:   url.LongURL,
		CreatorID: url.CreatorID,
	}
}

func CreateUrlInfoResponse(longUrl string, hitCount int64) UrlInfoResponse {
	return UrlInfoResponse{
		LongURL:  longUrl,
		HitCount: hitCount,
	}
}
