package controllers

import (
	"github.com/stretchr/testify/assert"
	"onepixel_backend/src/utils/applogger"
	_ "onepixel_backend/tests/providers"
	"testing"
)

var urlsController = CreateUrlsController()

func TestUrlsController_CreateSpecificShortUrl(t *testing.T) {
	user, _, _ := userController.Create("user136254@test.com", "123456")
	assert.NotNil(t, user)
	url, err := urlsController.CreateSpecificShortUrl("ax_bg", "https://example.com", user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.EqualValues(t, user.ID, url.CreatorID)
	applogger.Info("URL Created", "url ", url.ShortURL, "longUrl", url.LongURL)
}

func TestUrlsController_CreateRandomShortUrl(t *testing.T) {
	user, _, _ := userController.Create("user24987@test.com", "123456")

	assert.NotNil(t, user)
	url, err := urlsController.CreateRandomShortUrl("https://google.com", user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.EqualValues(t, user.ID, url.CreatorID)
	applogger.Info("URL Created", "url ", url.ShortURL, "longUrl", url.LongURL)
}
