package models

import "gorm.io/gorm"

// Url db entity
type Url struct {
	gorm.Model
	ID         uint64   `gorm:"primaryKey;autoIncrement:false"`
	ShortURL   string   `gorm:"unique,not null,size:10"`
	LongURL    string   `gorm:"not null"`
	CreatorID  uint64   `gorm:"not null"`
	Creator    User     `gorm:"foreignKey:CreatorID"`
	UrlGroupID uint64   `gorm:"primaryKey;not null, default:0"`
	UrlGroup   UrlGroup `gorm:"foreignKey:UrlGroupID"`
}

func (Url) TableName() string {
	return "urls"
}
