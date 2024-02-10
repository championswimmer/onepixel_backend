package models

import (
	"github.com/google/uuid"
	"time"
)

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
	// ip address
	IPAddress string `gorm:"type:string"`
	// referer
	Referer string `gorm:"type:string"`
}

func (EventRedirect) TableName() string {
	return "events_redirect"
}
