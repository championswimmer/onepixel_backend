package models

import (
	"github.com/google/uuid"
	"time"
)

// EventRedirect event db entity
type EventRedirect struct {
	CreatedAt time.Time
	//UUID key
	ID uuid.UUID `gorm:"primaryKey;type:UUID;default:generateUUIDv4()"`

	// short url
	ShortURL   string `gorm:"size:21"` // <group>/<shortcode>
	ShortUrlID uint64 `gorm:"bigint,not null"`
	UrlGroupID uint64 `gorm:"bigint,not null"`
	CreatorID  uint64 `gorm:"bigint,not null"`
	// user agent
	UserAgent string `gorm:"size:200"`
	// ip address
	IPAddress string `gorm:"size:20"`
	// referer
	Referer string `gorm:"size:255"`
}

func (EventRedirect) TableName() string {
	return "events_redirect"
}
