package controllers

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"onepixel_backend/src/db"
	"onepixel_backend/src/utils/applogger"
	"testing"
)

var urlsController = CreateUrlsController(lo.Must(db.GetTestDB()))

func TestUrlsController_CreateSpecificShortUrl(t *testing.T) {
	user, _, _ := userController.Create("user7353@test.com", "123456")
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
