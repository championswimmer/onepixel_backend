package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils"
)

// the current max length of the short url
// when we generate URLs randomly
// can be increase future if we run out of this space
//   - 6: 64^6 = 68,719,476,736
const _currentMaxUrlLength = 6

var _randMax = int(math.Pow(64, _currentMaxUrlLength))

type UrlsController struct {
	// db
	db *gorm.DB
}

type UrlError struct {
	status  int
	message string
}

func (e *UrlError) Error() string {
	return e.message
}

func (e *UrlError) ErrorDetails() (int, string) {
	return e.status, e.message
}

var (
	UrlNotFound = &UrlError{
		status:  404,
		message: "URL not found",
	}
	UrlExistsError = &UrlError{
		status:  fiber.ErrConflict.Code,
		message: "URL already exists",
	}
	UrlForbiddenError = &UrlError{
		status:  fiber.ErrForbidden.Code,
		message: "this shortURL is not allowed to be created",
	}
)

func CreateUrlsController(db *gorm.DB) *UrlsController {
	return &UrlsController{
		db: db,
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
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return nil, UrlExistsError
		} else {
			return nil, res.Error
		}
	}

	return
}

func (c *UrlsController) CreateRandomShortUrl(longUrl string, userId uint64) (url *models.Url, err error) {
	newShortCode := uint64(rand.Intn(_randMax))
	newShortUrl := lo.Must(utils.Radix64Encode(newShortCode))

	url = &models.Url{
		ID:         newShortCode,
		ShortURL:   newShortUrl,
		LongURL:    longUrl,
		CreatorID:  userId,
		UrlGroupID: nil,
	}
	res := c.db.Create(url)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			// Having a collision is very unlikely, but if it happens
			// we just try again (TODO: we should have a limit of retries)
			return c.CreateRandomShortUrl(longUrl, userId)
		}
		return nil, res.Error
	}

	return
}
