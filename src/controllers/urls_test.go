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

}
