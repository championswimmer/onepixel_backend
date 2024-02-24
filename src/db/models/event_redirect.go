package models

import (
	"github.com/google/uuid"
	"time"
)

type GeoIpData struct {
	LocationCity    string `gorm:"type:LowCardinality(String)"`
	LocationRegion  string `gorm:"type:LowCardinality(String)"`
	LocationCountry string `gorm:"type:LowCardinality(String)"`
}

// EventRedirect event db entity
type EventRedirect struct {
	CreatedAt time.Time
	//UUID key
	ID uuid.UUID `gorm:"primaryKey;"`

	// short url
	ShortURL   string `gorm:"type:string"` // <group>/<shortcode>
	ShortUrlID uint64 `gorm:"bigint,not null"`
	UrlGroupID uint64 `gorm:"bigint,not null"`
	CreatorID  uint64 `gorm:"bigint,not null"`
	// user agent
	UserAgent string `gorm:"type:string"`
	// ip address // TODO: have IPV4 and IPV6 columns
	IPAddress string `gorm:"type:string"`
	GeoIpData
	// referer
	Referer string `gorm:"type:string"`
}

// TableName returns the table name for the event redirect
func (EventRedirect) TableName() string {
	return "events_redirect"
}

// ----- row models for query views -----

type EventRedirectCountView struct {
	Redirects uint64 `gorm:"type:bigint"`
	ShortURL  string `gorm:"type:string"`
}
