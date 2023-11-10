package controllers

import (
	"errors"
	"gorm.io/gorm"
	"onepixel_backend/src/helpers"
	"onepixel_backend/src/models"
)

type UrlController struct {
	// db
	db *gorm.DB
}

func NewUrlController(db *gorm.DB) *UrlController {
	return &UrlController{
		db: db,
	}
}

// Create new url
func (c *UrlController) Create(longUrl string, groupId uint, userId uint) (*models.Url, error) {
	url := &models.Url{
		LongURL:   longUrl,
		GroupID:   groupId,
		CreatorID: userId,
	}

	shortUrl, err := helpers.GenerateShortUrl(6)
	if err != nil {
		return nil, err
	}
	url.ShortURL = shortUrl
	res := c.db.Create(url)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return nil, res.Error
		}
	}
	if res.Error != nil {
		return nil, res.Error
	}
	err = c.db.Preload("Creator").First(url, url.ID).Error
	if err != nil {
		return nil, err
	}
	return url, nil
}
