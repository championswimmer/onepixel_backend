package controllers

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
	"math/rand"
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

func (c *UrlsController) generateNewRandomShortUrl() string {
	newRandomUrlCode := uint64(rand.Intn(utils.MaxSafeNumber))
	var existingUrl models.Url
	c.db.Find(&existingUrl, models.Url{ID: newRandomUrlCode})
	if existingUrl.ID == newRandomUrlCode {
		return c.generateNewRandomShortUrl()
	} else {
		return lo.Must(utils.Radix64Encode(newRandomUrlCode))
	}
}

func (c *UrlsController) CreateSpecificShortUrl(shortUrl string, longUrl string, userId uint64) (url *models.Url, err error) {
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

func (c *UrlsController) CreateRandomShortUrl(longUrl string, userId uint64) {

}
