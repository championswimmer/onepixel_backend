package controllers

import (
	"onepixel_backend/src/server/validators"
	"onepixel_backend/src/utils/applogger"
	_ "onepixel_backend/tests/providers"
	"testing"

	"github.com/stretchr/testify/assert"
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

	t.Run("CreateEmptySpecificShortUrl", func(t *testing.T) {
		url, err := urlsController.CreateSpecificShortUrl("", "https://stackoverflow.com/", user.ID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, validators.ShortcodeEmptyError.Error())
		assert.Nil(t, url)
		applogger.Error("Could not create the URL: shortcode is empty")

	})

	t.Run("CreateMaxLengthExceedingSpecificShortUrl", func(t *testing.T) {
		url, err := urlsController.CreateSpecificShortUrl("mypinkreddit", "https://www.reddit.com/", user.ID)
		assert.NotNil(t, err)
		assert.EqualError(t, err, validators.ShortcodeTooLongError.Error())
		assert.Nil(t, url)
		applogger.Error("Could not create the URL: shortcode exceeds the maximum allowed length of 10 characters")

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
