package models

import "gorm.io/gorm"

// EventRedirect event db entity
type EventRedirect struct {
	gorm.Model
	// UUID key
	ID string `gorm:"primaryKey;autoIncrement:false"`
	// short url
	ShortURL  string `gorm:"not null,size:10"`
	LongURL   string `gorm:"not null"`
	CreatorID uint64 `gorm:"not null"`
	// user agent
	UserAgent string `gorm:"nullable"`
	// ip address
	IPAddress string `gorm:"nullable"`
	// referer
	Referer string `gorm:"nullable"`
}

func (EventRedirect) TableName() string {
	return "events_redirect"
}
