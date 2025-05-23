package controllers

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils/applogger"
	_ "onepixel_backend/tests/providers"
	"testing"
	"time"
)

var urlsController = CreateUrlsController()
var eventsController = CreateEventsController()

func TestUrlsController(t *testing.T) {
	user, _, _ := userController.Create("user14612@test.com", "123456")
	assert.NotNil(t, user)

	t.Run("CreateSpecificShortUrl", func(t *testing.T) {
		url, err := urlsController.CreateSpecificShortUrl("ax_bg", "https://example.com", user.ID)
		assert.Nil(t, err)
		assert.NotNil(t, url)
		assert.EqualValues(t, user.ID, url.CreatorID)
		applogger.Info("URL Created", "url ", url.ShortURL, "longUrl", url.LongURL)

	})

	t.Run("CreateRandomShortUrl", func(t *testing.T) {
		url, err := urlsController.CreateRandomShortUrl("https://google.com", user.ID)
		assert.Nil(t, err)
		assert.NotNil(t, url)
		assert.EqualValues(t, user.ID, url.CreatorID)
		applogger.Info("URL Created", "url ", url.ShortURL, "longUrl", url.LongURL)

	})

	t.Run("CreateUrlGroup", func(t *testing.T) {
		urlGroup, err := urlsController.CreateUrlGroup("grp", user.ID)
		assert.Nil(t, err)
		assert.NotNil(t, urlGroup)
	})

	t.Run("GetUrlInfo", func(t *testing.T) {
		// Add a new URL
		url, err := urlsController.CreateSpecificShortUrl("test123", "https://example.com", user.ID)
		assert.Nil(t, err)
		assert.NotNil(t, url)

		// Fetch the URL info
		longUrl, hitCount, err := urlsController.GetUrlInfo("test123")
		assert.Nil(t, err)
		assert.Equal(t, "https://example.com", longUrl)
		assert.Equal(t, int64(0), hitCount)

		// Simulate a hit
		eventsController.LogRedirectAsync(&EventRedirectData{
			ShortUrlID: url.ID,
			UrlGroupID: url.UrlGroupID,
			ShortURL:   url.ShortURL,
			CreatorID:  url.CreatorID,
			IPAddress:  "127.0.0.1",
			UserAgent:  "test-agent",
			Referer:    "test-referer",
		})

		time.Sleep(200 * time.Millisecond)

		// Fetch the URL info again
		longUrl, hitCount, err = urlsController.GetUrlInfo("test123")
		assert.Nil(t, err)
		assert.Equal(t, "https://example.com", longUrl)
		assert.Equal(t, int64(1), hitCount)
	})

	t.Run("GetUrls", func(t *testing.T) {
		// Add a new URL
		url, err := urlsController.CreateSpecificShortUrl("test456", "https://example.com", user.ID)
		assert.Nil(t, err)
		assert.NotNil(t, url)

		// Fetch URLs by user ID
		urls, err := urlsController.GetUrls(&user.ID)
		assert.Nil(t, err)
		assert.NotNil(t, urls)
		assert.Greater(t, len(urls), 0)
		assert.Equal(t, user.ID, urls[0].CreatorID)
	})

	t.Run("GetAllUrlsByAdmin", func(t *testing.T) {
		// Add a new URL
		url, err := urlsController.CreateSpecificShortUrl("test777", "https://example.com", user.ID)
		assert.Nil(t, err)
		assert.NotNil(t, url)

		// Fetch all URLs
		var u *uint64
		urls, err := urlsController.GetUrls(u)
		assert.Nil(t, err)
		assert.Equal(t, len(urls), 1)
		_, found := lo.Find(urls, func(item models.Url) bool {
			return item.ShortURL == "test777"
		})
		assert.True(t, found)
	})
}
