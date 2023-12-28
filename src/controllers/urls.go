package controllers

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
	"onepixel_backend/src/models"
	"onepixel_backend/src/utils"
)

type UrlsController struct {
	// db
	db *gorm.DB
}

func CreateUrlsController(db *gorm.DB) *UrlsController {
	return &UrlsController{
		db: db,
	}
}

func (c *UrlsController) CreateShortUrl(shortUrl string, longUrl string, userId uint) (url *models.Url, err error) {
	url = &models.Url{
		ID:         lo.Must(utils.Radix64Decode(shortUrl)),
		ShortURL:   shortUrl,
		LongURL:    longUrl,
		CreatorID:  userId,
		UrlGroupID: nil,
	}

	res := c.db.Create(url)
	if res.Error != nil {
		return nil, res.Error
	}

	return
}
