package controllers

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"onepixel_backend/src/models"
	"onepixel_backend/src/utils"
)

const CURRENT_SOFTMAX_URL_LENGTH = 6

var RANDMAX = int(math.Pow(64, CURRENT_SOFTMAX_URL_LENGTH))

type UrlsController struct {
	// db
	db *gorm.DB
}

func CreateUrlsController(db *gorm.DB) *UrlsController {
	return &UrlsController{
		db: db,
	}
}

func (c *UrlsController) generateNewRandomShortUrl() (string, uint64) {
	// TODO: optimise in future: why check before trying to create itself -
	// 		if it fails with unique constraint then only recreate
	newRandomUrlCode := uint64(rand.Intn(RANDMAX))
	var existingUrl models.Url
	c.db.Find(&existingUrl, models.Url{ID: newRandomUrlCode})
	if existingUrl.ID == newRandomUrlCode {
		return c.generateNewRandomShortUrl()
	} else {
		return lo.Must(utils.Radix64Encode(newRandomUrlCode)), newRandomUrlCode
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

func (c *UrlsController) CreateRandomShortUrl(longUrl string, userId uint64) (url *models.Url, err error) {
	newShortUrl, newShortCode := c.generateNewRandomShortUrl()

	url = &models.Url{
		ID:         newShortCode,
		ShortURL:   newShortUrl,
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
