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
const _maxRandomShortUrlRetries = 10

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
	UrlGroupExistsError = &UrlError{
		status:  fiber.ErrConflict.Code,
		message: "URL group already exists",
	}
	UrlForbiddenError = &UrlError{
		status:  fiber.ErrForbidden.Code,
		message: "this shortURL is not allowed to be created",
	}
	UrlGroupForbiddenError = &UrlError{
		status:  fiber.ErrForbidden.Code,
		message: "URL group does not belong to the user",
	}
	UrlGroupNotFound = &UrlError{
		status:  fiber.StatusNotFound,
		message: "URL group not found",
	}
	RandomShortUrlExhaustedError = &UrlError{
		status:  fiber.StatusInternalServerError,
		message: "failed to generate a unique short URL",
	}
)

var initDefaultUrlGroupOnce sync.Once

func (c *UrlsController) initDefaultUrlGroup() {
	defaultUrlGroup := &models.UrlGroup{
		ID:        _defaultUrlGroupId,
		ShortPath: lo.Must(utils.Radix64Encode(_defaultUrlGroupId)), // "0",
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
	return c.CreateSpecificShortUrlInGroup(shortUrl, longUrl, userId, _defaultUrlGroupId)
}

func (c *UrlsController) CreateSpecificShortUrlInGroup(shortUrl string, longUrl string, userId uint64, urlGroupID uint64) (url *models.Url, err error) {
	if err := c.ensureUrlGroupOwnership(urlGroupID, userId); err != nil {
		return nil, err
	}

	url = &models.Url{
		ID:         lo.Must(utils.Radix64Decode(shortUrl)),
		ShortURL:   shortUrl,
		LongURL:    longUrl,
		CreatorID:  userId,
		UrlGroupID: urlGroupID,
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
	return c.CreateRandomShortUrlInGroup(longUrl, userId, _defaultUrlGroupId)
}

func (c *UrlsController) CreateRandomShortUrlInGroup(longUrl string, userId uint64, urlGroupID uint64) (url *models.Url, err error) {
	if err := c.ensureUrlGroupOwnership(urlGroupID, userId); err != nil {
		return nil, err
	}

	for attempt := 0; attempt < _maxRandomShortUrlRetries; attempt++ {
		newShortCode := uint64(rand.Intn(_randMax))
		newShortUrl := lo.Must(utils.Radix64Encode(newShortCode))

		url = &models.Url{
			ID:         newShortCode,
			ShortURL:   newShortUrl,
			LongURL:    longUrl,
			CreatorID:  userId,
			UrlGroupID: urlGroupID,
		}

		res := c.db.Create(url)
		if res.Error == nil {
			return url, nil
		}
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			continue
		}
		return nil, res.Error
	}

	return nil, RandomShortUrlExhaustedError
}

func (c *UrlsController) GetUrlWithShortCode(shortcode string) (url *models.Url, err error) {
	return c.GetUrlWithShortCodeInGroup(shortcode, _defaultUrlGroupId)
}

func (c *UrlsController) GetUrlWithShortCodeInGroup(shortcode string, urlGroupID uint64) (url *models.Url, err error) {
	url = &models.Url{}
	res := c.db.Where("url_group_id = ? AND short_url = ?", urlGroupID, shortcode).First(url)
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
		ShortPath: groupName,
		ID:        lo.Must(utils.Radix64Decode(groupName)),
		CreatorID: userId,
	}

	res := c.db.Create(urlGroup)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return nil, UrlGroupExistsError
		}
		return nil, res.Error
	}
	return urlGroup, nil

}

func (c *UrlsController) GetUrlGroupByShortPath(groupName string) (urlGroup *models.UrlGroup, err error) {
	urlGroup = &models.UrlGroup{}
	res := c.db.Where("name = ?", groupName).First(urlGroup)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, UrlGroupNotFound
		}
		return nil, res.Error
	}
	return urlGroup, nil
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
	query := c.db.Preload("UrlGroup")
	if userId != nil {
		query = query.Where("creator_id=?", *userId)
	}
	res := query.Find(&urls)
	if res.Error != nil {
		return nil, res.Error
	}
	return urls, nil
}

func (c *UrlsController) ensureUrlGroupOwnership(urlGroupID uint64, userID uint64) error {
	if urlGroupID == _defaultUrlGroupId {
		return nil
	}

	urlGroup := &models.UrlGroup{}
	res := c.db.First(urlGroup, urlGroupID)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return UrlGroupNotFound
		}
		return res.Error
	}
	if urlGroup.CreatorID != userID {
		return UrlGroupForbiddenError
	}
	return nil
}
