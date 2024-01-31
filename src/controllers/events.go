package controllers

import (
	"github.com/samber/lo"
	"gorm.io/gorm"
	"onepixel_backend/src/db"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils/applogger"
	"sync"
)

type EventsController struct {
	// event logging eventDb (not the main app eventDb)
	eventDb *gorm.DB
}

func CreateEventsController() *EventsController {
	eventsDB := lo.Must(db.GetEventsDB())
	return &EventsController{
		eventDb: eventsDB,
	}
}

type EventRedirectDTO struct {
	ShortUrlID uint64
	UrlGroupID uint64
	ShortURL   string
	CreatorID  uint64
	IPAddress  string
	UserAgent  string
	Referer    string
}

var wg = sync.WaitGroup{}

func (c *EventsController) LogRedirectAsync(redirData *EventRedirectDTO) {
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				applogger.Error("LogRedirectAsync panic: ", r)
			}
			wg.Done()
		}()

		event := &models.EventRedirect{
			ShortURL:   redirData.ShortURL,
			ShortUrlID: redirData.ShortUrlID,
			UrlGroupID: redirData.UrlGroupID,
			CreatorID:  redirData.CreatorID,
			UserAgent:  redirData.UserAgent,
			IPAddress:  redirData.IPAddress,
			Referer:    redirData.Referer,
		}
		c.eventDb.Create(event)

	}()
}
