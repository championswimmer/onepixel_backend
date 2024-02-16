package controllers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"net"
	"onepixel_backend/src/db"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/utils/applogger"
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
		ip := net.ParseIP(redirData.IPAddress)

		if ip != nil {
			city, err := c.geoipDB.City(ip)
			if err == nil {
				if city.Country.Names["en"] != "" {
					event.LocationCountry = fmt.Sprintf("%s (%s)", city.Country.Names["en"], city.Country.IsoCode)
				}

				if city.Subdivisions[0].Names["en"] != "" {
					event.LocationRegion = fmt.Sprintf("%s (%s)", city.Subdivisions[0].Names["en"], city.Subdivisions[0].IsoCode)
				}

				if city.City.Names["en"] != "" {
					event.LocationCity = city.City.Names["en"]
				}
			}
		}
		lo.Try(func() error {
			tx := c.eventDb.Create(event)
			return tx.Error
		})
		return event.ID
	})

}

func (c *EventsController) GetRedirectsCountForUserId(userId uint64) ([]models.EventRedirectCountView, error) {
	rows, err := c.eventDb.Model(&models.EventRedirect{}).
		Select("count(id) as redirects, short_url").
		Where("creator_id = ?", userId).
		Group("short_url").
		Rows()

	if err != nil {
		applogger.Error("GetRedirectsCountForUserId: ", err)
		return nil, err
	}
	data := make([]models.EventRedirectCountView, 0)
	for rows.Next() {
		var d models.EventRedirectCountView
		lo.Must0(c.eventDb.Model(&models.EventRedirect{}).ScanRows(rows, &d))
		data = append(data, d)
	}
	return data, nil
}
