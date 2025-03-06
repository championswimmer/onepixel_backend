package controllers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"onepixel_backend/src/db"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils/applogger"
	"onepixel_backend/src/utils/clientinfo"
)

type EventsController struct {
	// event logging eventDb (not the main app eventDb)
	eventDb *gorm.DB
	geoipDB *geoip2.Reader
}

func CreateEventsController() *EventsController {
	return &EventsController{
		eventDb: db.GetEventsDB(),
		geoipDB: db.GetGeoIPDB(),
	}
}

type EventRedirectData struct {
	ShortUrlID uint64
	UrlGroupID uint64
	ShortURL   string
	CreatorID  uint64
	IPAddress  string
	UserAgent  string
	Referer    string
}

// synced is a synchronized locker to access user-specific views
var synced = lo.Synchronize()

func (c *EventsController) getViewForUser(userId uint64) *gorm.DB {
	var evdb *gorm.DB
	synced.Do(func() {
		viewName := (&models.EventRedirect{}).TableName() + "_" + fmt.Sprint(userId)
		viewOption := gorm.ViewOption{
			Replace: true,
			Query:   c.eventDb.Model(&models.EventRedirect{}).Where("creator_id = ?", userId),
		}
		lo.Try(func() error {
			err := c.eventDb.Migrator().CreateView(viewName, viewOption)
			if err != nil {
				applogger.Error("getViewForUser: failed to create view for user id: ", userId, err)
			}
			return err
		})
		evdb = c.eventDb.Table(viewName)
	})
	return evdb
}

func (c *EventsController) LogRedirectAsync(redirData *EventRedirectData) {
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
		lo.TryCatchWithErrorValue(func() error {
			geoIpData, err := clientinfo.GetGeoIpDataFromIP(c.geoipDB, redirData.IPAddress)
			event.GeoIpData = *geoIpData
			return err
		}, func(e any) {
			applogger.Warn("LogRedirectAsync: failed to get geoip data: ", e)
		})

		lo.Try(func() error {
			tx := c.eventDb.Create(event)
			if tx.Error != nil {
				applogger.Error("LogRedirectAsync: failed to create event redirect: ", tx.Error)
			}
			return tx.Error
		})
		return event.ID
	})

}

func (c *EventsController) GetRedirectsCountForUserId(userId uint64) ([]models.EventRedirectCountView, error) {
	view := c.getViewForUser(userId)
	rows, err := view.Model(&models.EventRedirect{}).
		Select("count(id) as redirects, short_url").
		Group("short_url").
		Rows()

	if err != nil {
		applogger.Error("GetRedirectsCountForUserId: ", err)
		return nil, err
	}
	data := make([]models.EventRedirectCountView, 0)
	for rows.Next() {
		var d models.EventRedirectCountView
		lo.Must0(view.Model(&models.EventRedirect{}).ScanRows(rows, &d))
		data = append(data, d)
	}
	return data, nil
}

/*
```clickhouse
select count(id)                                 as redirects,
       location_city || ', ' || location_country as city,
       short_url
from events_redirect
where created_at >= '2024-02-18'
  and created_at <= '2024-02-19'
group by short_url,
         location_city,
         location_country
order by redirects desc
```
*/
