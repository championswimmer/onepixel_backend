package controllers

import (
	"errors"
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

	t.Run("CreateUrlGroupDuplicateFail", func(t *testing.T) {
		_, err := urlsController.CreateUrlGroup("dupgrp", user.ID)
		assert.Nil(t, err)

		_, err = urlsController.CreateUrlGroup("dupgrp", user.ID)
		assert.ErrorIs(t, err, UrlGroupExistsError)
	})

	t.Run("GroupedUrlsRequireOwnership", func(t *testing.T) {
		otherUser, _, err := userController.Create("user24680@test.com", "123456")
		assert.Nil(t, err)

		urlGroup, err := urlsController.CreateUrlGroup("grpown", user.ID)
		assert.Nil(t, err)

		groupedUrl, err := urlsController.CreateSpecificShortUrlInGroup("ownabc", "https://owner.example.com", user.ID, urlGroup.ID)
		assert.Nil(t, err)
		assert.EqualValues(t, urlGroup.ID, groupedUrl.UrlGroupID)

		_, err = urlsController.CreateRandomShortUrlInGroup("https://forbidden.example.com", otherUser.ID, urlGroup.ID)
		assert.ErrorIs(t, err, UrlGroupForbiddenError)
	})

	t.Run("GroupedAndUngroupedLookupsAreScoped", func(t *testing.T) {
		urlGroup, err := urlsController.CreateUrlGroup("grpdup", user.ID)
		assert.Nil(t, err)

		groupedUrl, err := urlsController.CreateSpecificShortUrlInGroup("same01", "https://group.example.com", user.ID, urlGroup.ID)
		assert.Nil(t, err)
		assert.EqualValues(t, urlGroup.ID, groupedUrl.UrlGroupID)

		ungroupedUrl, err := urlsController.CreateSpecificShortUrl("same01", "https://default.example.com", user.ID)
		assert.Nil(t, err)
		assert.EqualValues(t, uint64(0), ungroupedUrl.UrlGroupID)

		defaultLookup, err := urlsController.GetUrlWithShortCode("same01")
		assert.Nil(t, err)
		assert.Equal(t, "https://default.example.com", defaultLookup.LongURL)
		assert.EqualValues(t, uint64(0), defaultLookup.UrlGroupID)

		groupedLookup, err := urlsController.GetUrlWithShortCodeInGroup("same01", urlGroup.ID)
		assert.Nil(t, err)
		assert.Equal(t, "https://group.example.com", groupedLookup.LongURL)
		assert.EqualValues(t, urlGroup.ID, groupedLookup.UrlGroupID)

		lookupGroup, err := urlsController.GetUrlGroupByShortPath("grpdup")
		assert.Nil(t, err)
		assert.EqualValues(t, urlGroup.ID, lookupGroup.ID)
	})

	t.Run("MissingUrlGroupReturnsTypedError", func(t *testing.T) {
		_, err := urlsController.GetUrlGroupByShortPath("missinggrp")
		assert.True(t, errors.Is(err, UrlGroupNotFound))
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
		assert.Greater(t, len(urls), 0)
		_, found := lo.Find(urls, func(item models.Url) bool {
			return item.ShortURL == "test777"
		})
		assert.True(t, found)
		groupedUrl, ok := lo.Find(urls, func(item models.Url) bool {
			return item.UrlGroupID != 0
		})
		assert.True(t, ok)
		assert.NotEmpty(t, groupedUrl.UrlGroup.ShortPath)
	})
}
