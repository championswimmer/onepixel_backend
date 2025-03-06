package controllers

import (
	"github.com/stretchr/testify/assert"
	"onepixel_backend/src/utils/applogger"
	_ "onepixel_backend/tests/providers"
	"testing"
)

var urlsController = CreateUrlsController()

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

		// Fetch the URL info again
		longUrl, hitCount, err = urlsController.GetUrlInfo("test123")
		assert.Nil(t, err)
		assert.Equal(t, "https://example.com", longUrl)
		assert.Equal(t, int64(1), hitCount)
	})
}
