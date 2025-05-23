package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"onepixel_backend/src/db"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils"
	"onepixel_backend/src/utils/applogger"
	"sync"
)

// the current max length of the short url
// when we generate URLs randomly
// can be increase future if we run out of this space
//   - 6: 64^6 = 68,719,476,736
const _currentMaxUrlLength = 6
const _defaultUrlGroupId = 0

var _randMax = int(math.Pow(64, _currentMaxUrlLength))

type UrlsController struct {
	// db
	db               *gorm.DB
	eventsController *EventsController
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

var initDefaultUrlGroupOnce sync.Once

func (c *UrlsController) initDefaultUrlGroup() {
	defaultUrlGroup := &models.UrlGroup{
		ID:        _defaultUrlGroupId,
		Name:      lo.Must(utils.Radix64Encode(_defaultUrlGroupId)), // "0",
		CreatorID: 0,
	}

	res := c.db.Save([]models.UrlGroup{*(defaultUrlGroup)})
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			applogger.Warn("Default url group already exists")
			return
		}
		applogger.Error("Failed to create default url group")
		applogger.Panic(res.Error)
	} else {
		applogger.Info("Default url group created")
	}
}

func CreateUrlsController() *UrlsController {
	appDb := db.GetAppDB()
	ctrl := &UrlsController{
		db:               appDb,
		eventsController: CreateEventsController(),
	}
	initDefaultUrlGroupOnce.Do(ctrl.initDefaultUrlGroup)
	return ctrl
}

func (c *UrlsController) CreateSpecificShortUrl(shortUrl string, longUrl string, userId uint64) (url *models.Url, err error) {
	url = &models.Url{
		ID:         lo.Must(utils.Radix64Decode(shortUrl)),
		ShortURL:   shortUrl,
		LongURL:    longUrl,
		CreatorID:  userId,
		UrlGroupID: 0,
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
		UrlGroupID: 0,
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

func (c *UrlsController) GetUrlWithShortCode(shortcode string) (url *models.Url, err error) {
	url = &models.Url{}
	id := lo.Must(utils.Radix64Decode(shortcode))
	res := c.db.First(url, id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, UrlNotFound
		}
		return nil, res.Error
	}

	return
}

func (c *UrlsController) CreateUrlGroup(groupName string, userId uint64) (urlGroup *models.UrlGroup, err error) {
	urlGroup = &models.UrlGroup{
		Name:      groupName,
		ID:        lo.Must(utils.Radix64Decode(groupName)),
		CreatorID: userId,
	}

	res := c.db.Create(urlGroup)
	if res.Error != nil {

	}
	return

}

func (c *UrlsController) GetUrlInfo(shortcode string) (longUrl string, hitCount int64, err error) {
	url, err := c.GetUrlWithShortCode(shortcode)
	if err != nil {
		return "", 0, err
	}

	hitCount, err = c.eventsController.GetRedirectsCountForShortCode(shortcode)
	if err != nil {
		return "", 0, err
	}

	return url.LongURL, hitCount, nil
}

func (c *UrlsController) GetUrls(userId *uint64) ([]models.Url, error) {
	var urls []models.Url
	query := c.db
	if userId != nil {
		query = query.Where("creator_id=?", *userId)
	}
	res := query.Find(&urls)
	if res.Error != nil {
		return nil, res.Error
	}
	return urls, nil
}
