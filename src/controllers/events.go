package controllers

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"onepixel_backend/src/db"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils/applogger"
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

func (c *EventsController) LogRedirectAsync(redirData *EventRedirectDTO) {

	lo.Async(func() uuid.UUID {
		event := &models.EventRedirect{
			ID:         uuid.New(),
			ShortURL:   redirData.ShortURL,
			ShortUrlID: redirData.ShortUrlID,
			UrlGroupID: redirData.UrlGroupID,
			CreatorID:  redirData.CreatorID,
			UserAgent:  redirData.UserAgent,
			IPAddress:  redirData.IPAddress,
			Referer:    redirData.Referer,
		}
		lo.Try(func() error {
			tx := c.eventDb.Create(event)
			return tx.Error
		})
		return event.ID
	})

}

func (c *EventsController) GetRedirectsCountForUserId(userId string) []models.EventRedirectCountView {
	rows, err := c.eventDb.Model(&models.EventRedirect{}).
		Select("count(id) as redirects, short_url").
		Group("short_url").
		Rows()

	if err != nil {
		applogger.Panic("GetRedirectsCountForUserId: ", err)
	}
	data := make([]models.EventRedirectCountView, 0)
	for rows.Next() {
		var d models.EventRedirectCountView
		lo.Must0(c.eventDb.Model(&models.EventRedirect{}).ScanRows(rows, &d))
		data = append(data, d)
	}
	return data
}
