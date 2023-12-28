package controllers

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"onepixel_backend/src/db"
	"testing"
)

var urlsController = CreateUrlsController(lo.Must(db.GetTestDB()))

func TestUrlsController_CreateShortUrl(t *testing.T) {
	user, _, _ := userController.Create("user7353@test.com", "123456")
	assert.NotNil(t, user)
	url, err := urlsController.CreateSpecificShortUrl("ax_bg", "https://example.com", user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, url)
	assert.EqualValues(t, user.ID, url.CreatorID)

}
